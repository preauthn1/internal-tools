// 模型服務：只包 Apple 本機模型的 HTTP 服務，給 Go 版呼叫（go/internal/model 的 Apple）
// 編譯與執行：swift build -c release && .build/release/model-server（環境變數見 README「模型層」）
//   POST /respond {name, instructions, text, schema}：成功回 {values, usage}；本機模型失敗或逾時回 502 {error}，
//     提示詞超過本機模型上下文時 502 多帶 reason: "context_exceeded"，被安全過濾擋下時帶 reason: "guardrail"
//     usage {input_tokens, output_tokens}：這次用掉的 token 數（Go 端記進 AI 費用頁）；算不出來時沒有這個欄位
//     排隊上限：已經有 MODEL_SERVER_QUEUE_MAX（預設 4）個在等時，新的請求直接回 503，不排隊
//     X-Start-By-Ms（Unix 毫秒）：排隊到這個時間還沒輪到就不執行，回 504（請求方已經等不到結果，不要白佔本機模型）
//   GET  /version：{model, os}，Go 端放進 AI 回答快取的 key
// 只聽本機（MODEL_SERVER_HOST，預設 0.0.0.0），另外檢查連線來源：只接受 loopback 和 MODEL_SERVER_ALLOW 列的網段
// 呼叫經過 AIQueue：本機模型一次只跑一個，和 Swift 版的量測程式共用同一把檔案鎖。
// 上限和放棄只在這裡做，不改 BotCore 的 QueuedModel：量測程式本來就要「排再久也等到」
// log 寫到 stderr（不緩衝，launchd 導到 ~/Library/Logs/aichatbot/model-server.log），每行開頭是本地時間到秒；
// 回 502 時寫出是逾時還是本機模型丟出錯誤（型別和訊息），以及推論花了多久
import Foundation
import FoundationModels
import Hummingbird
import NIOCore
import BotCore

// 時間格式和 Go 版 log 對得起來：2026-10-12 14:03:05（本地時間）
let timestampStyle = Date.ISO8601FormatStyle(dateSeparator: .dash, dateTimeSeparator: .space, timeZone: .current)
    .year().month().day().time(includingFractionalSeconds: false)
@Sendable func log(_ s: String) {
    FileHandle.standardError.write(Data("\(Date().formatted(timestampStyle)) \(s)\n".utf8))
}

guard case .available = SystemLanguageModel.default.availability else {
    log("❌ 本機模型無法使用：\(SystemLanguageModel.default.availability)"); exit(1)
}

let vars = ProcessInfo.processInfo.environment
let host = vars["MODEL_SERVER_HOST"] ?? "0.0.0.0"
guard let port = Int(vars["MODEL_SERVER_PORT"] ?? "8765") else {
    log("❌ 環境變數 MODEL_SERVER_PORT 不正確：要是數字"); exit(1)
}
// MODEL_SERVER_ALLOW：逗號分隔的 IPv4 網段（例如 0.0.0.0/24），loopback 一律允許
var allowed: [(network: UInt32, mask: UInt32)] = []
for item in (vars["MODEL_SERVER_ALLOW"] ?? "").split(separator: ",") {
    guard let cidr = parseCIDR(item.trimmingCharacters(in: .whitespaces)) else {
        log("❌ 環境變數 MODEL_SERVER_ALLOW 不正確：\(item)（要像 0.0.0.0/24，逗號分隔）"); exit(1)
    }
    allowed.append(cidr)
}
guard let queueMax = Int(vars["MODEL_SERVER_QUEUE_MAX"] ?? "4"), queueMax >= 0 else {
    log("❌ 環境變數 MODEL_SERVER_QUEUE_MAX 不正確：要是 0 以上的整數"); exit(1)
}
let apple = AppleModel()

// 正在跑的 1 個加上排隊的，最多 queueMax + 1 個；滿了就拒絕
actor Admission {
    let max: Int
    var inSystem = 0
    init(max: Int) { self.max = max }
    func enter() -> Bool {
        guard inSystem <= max else { return false }
        inSystem += 1
        return true
    }
    func leave() { inSystem -= 1 }
}
let admission = Admission(max: queueMax)

enum Outcome: Sendable { case expired, done(TimedResult<(reply: ModelReply, usage: ModelUsage?)>, seconds: TimeInterval) }

// 成功的回應：values 是各欄位的值，usage 是 token 數（沒有就不放）
struct RespondBody: Encodable {
    let values: [String: FieldValue]
    let usage: ModelUsage?
}

@Sendable func ipv4(_ s: String) -> UInt32? {
    var addr = in_addr()
    guard inet_pton(AF_INET, s, &addr) == 1 else { return nil }
    return UInt32(bigEndian: addr.s_addr)
}

func parseCIDR(_ s: String) -> (network: UInt32, mask: UInt32)? {
    let parts = s.split(separator: "/")
    guard parts.count == 2, let ip = ipv4(String(parts[0])), let bits = Int(parts[1]), (0...32).contains(bits) else { return nil }
    let mask: UInt32 = bits == 0 ? 0 : ~UInt32(0) << (32 - bits)
    return (ip & mask, mask)
}

@Sendable func isAllowed(_ address: SocketAddress?) -> Bool {
    guard let ip = address?.ipAddress else { return false }
    if ip == "::1" { return true }
    // IPv6 對應的 IPv4（::ffff:0.0.0.0）當 IPv4 看
    let v4 = ip.hasPrefix("::ffff:") ? String(ip.dropFirst(7)) : ip
    guard let n = ipv4(v4) else { return false }
    if n >> 24 == 127 { return true }
    return allowed.contains { n & $0.mask == $0.network }
}

struct ModelContext: RequestContext, RemoteAddressRequestContext {
    var coreContext: CoreRequestContextStorage
    let remoteAddress: SocketAddress?
    init(source: ApplicationRequestContextSource) {
        coreContext = .init(source: source)
        remoteAddress = source.channel.remoteAddress
    }
}

struct SourceCheck: RouterMiddleware {
    func handle(_ request: Request, context: ModelContext, next: (Request, ModelContext) async throws -> Response) async throws -> Response {
        guard isAllowed(context.remoteAddress) else {
            log("⚠️ 拒絕連線：\(context.remoteAddress?.ipAddress ?? "不明來源")")
            return Response(status: .forbidden)
        }
        return try await next(request, context)
    }
}

@Sendable func json(_ status: HTTPResponse.Status, _ object: Any) -> Response {
    let data = (try? JSONSerialization.data(withJSONObject: object)) ?? Data("{}".utf8)
    return Response(status: status, headers: [.contentType: "application/json"], body: .init(byteBuffer: ByteBuffer(bytes: data)))
}

let router = Router(context: ModelContext.self)
router.add(middleware: SourceCheck())
router.get("/version") { _, _ in
    json(.ok, ["model": AppleModel.version, "os": ProcessInfo.processInfo.operatingSystemVersionString])
}
router.post("/respond") { request, _ in
    let body = try await request.body.collect(upTo: 1 << 20)
    guard let object = try? JSONSerialization.jsonObject(with: Data(body.readableBytesView)) as? [String: Any],
          let name = object["name"] as? String, let instructions = object["instructions"] as? String,
          let text = object["text"] as? String, let schema = object["schema"], let fields = fieldsFromJSONSchema(schema) else {
        return json(.badRequest, ["error": "請求要有 name、instructions、text 和 schema（Go 版 model.Schema 的格式）"])
    }
    let startBy = request.headers[.init("X-Start-By-Ms")!].flatMap(Double.init).map { Date(timeIntervalSince1970: $0 / 1000) }
    guard await admission.enter() else {
        log("⚠️ 排隊已滿（\(queueMax) 個在等），回 503：\(name)")
        return json(.serviceUnavailable, ["error": "模型服務忙碌：已經有 \(queueMax) 個請求在排隊"])
    }
    let request = ModelRequest(name: name, text: text, instructions: instructions, fields: fields)
    let outcome: Outcome = await AIQueue.shared.run {
        if let startBy, Date() > startBy { return .expired }
        let start = Date()
        let result = await apple.respondDetailed(request)
        return .done(result, seconds: Date().timeIntervalSince(start))
    }
    await admission.leave()
    switch outcome {
    case .expired:
        log("⚠️ 排隊太久，請求方已經等不到結果，不執行，回 504：\(name)")
        return json(.gatewayTimeout, ["error": "排隊太久，沒有執行"])
    case .done(.timedOut, let seconds):
        log("⚠️ 本機模型逾時（超過 \(Int(AppleModel.timeout)) 秒），回 502：\(name)，花了 \(String(format: "%.1f", seconds)) 秒")
        return json(.badGateway, ["error": "本機模型沒有回答（失敗或逾時）"])
    case .done(.failed(contextExceededType, _), let seconds):
        // Go 端看 reason 分出是「太長」，不是一般的 AI 失敗（#65）
        log("⚠️ 提示詞超過本機模型上下文，回 502：\(name)，指示 \(instructions.count) 字，花了 \(String(format: "%.1f", seconds)) 秒")
        return json(.badGateway, ["error": "提示詞超過本機模型上下文", "reason": "context_exceeded"])
    case .done(.failed(guardrailType, _), let seconds):
        // Go 端看 reason 分出是「被安全過濾擋下」，不是一般的 AI 失敗（#68）
        log("⚠️ 被安全過濾擋下，回 502：\(name)，花了 \(String(format: "%.1f", seconds)) 秒")
        return json(.badGateway, ["error": "被安全過濾擋下", "reason": "guardrail"])
    case .done(.failed(let type, let message), let seconds):
        // 錯誤訊息可能有換行，壓成一行
        let oneLine = message.split(whereSeparator: \.isNewline).joined(separator: " ")
        log("⚠️ 本機模型錯誤，回 502：\(name)：\(type)：\(oneLine)，花了 \(String(format: "%.1f", seconds)) 秒")
        return json(.badGateway, ["error": "本機模型沒有回答（失敗或逾時）"])
    case .done(.value(let answer), _):
        let body = try JSONEncoder().encode(RespondBody(values: answer.reply.values, usage: answer.usage))
        return Response(status: .ok, headers: [.contentType: "application/json"], body: .init(byteBuffer: ByteBuffer(bytes: body)))
    }
}

let app = Application(
    router: router,
    configuration: .init(address: .hostname(host, port: port)),
    onServerRunning: { channel in
        log("模型服務已啟動：http://\(host):\(channel.localAddress?.port ?? port)（\(AppleModel.version)）")
    }
)
try await app.runService()
