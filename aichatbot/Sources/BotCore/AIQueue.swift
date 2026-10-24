import Foundation

// MARK: - AI 排隊：同一台機器一次只跑一個 AI 呼叫（本機模型同時跑兩個，品質會明顯變差）
// 程序內先排隊，再用檔案鎖（flock）和其他程序排隊。鎖檔：BOT_AI_LOCK，預設在使用者快取目錄
// BOT_AI_LOG=檔案路徑：每個 AI 呼叫附加一行「程序編號<TAB>開始<TAB>結束」（Unix 秒，到毫秒），給 eval/check_ai_log.py 檢查有沒有重疊

public final class AIQueue: @unchecked Sendable {
    let lockPath: String
    let logPath: String?
    private let state = NSLock()
    private var busy = false
    private var waiters: [CheckedContinuation<Void, Never>] = []
    private var waited = 0.0, calls = 0

    public init(lockPath: String, logPath: String? = nil) {
        self.lockPath = lockPath; self.logPath = logPath
        try? FileManager.default.createDirectory(at: URL(fileURLWithPath: lockPath).deletingLastPathComponent(), withIntermediateDirectories: true)
    }

    public static let shared = AIQueue(
        lockPath: env["BOT_AI_LOCK"] ?? FileManager.default.urls(for: .cachesDirectory, in: .userDomainMask)[0]
            .appendingPathComponent("aichatbot/ai.lock").path,
        logPath: env["BOT_AI_LOG"])

    // 等鎖的總時間（秒）和 AI 呼叫次數（量測結果用）
    public var lockWaitSeconds: Double { state.withLock { waited } }
    public var callCount: Int { state.withLock { calls } }

    // 拿到兩層鎖才執行 op；逾時算在 op 裡面，所以等鎖不會讓呼叫變成逾時。
    // 逾時時 op 先回傳 nil、鎖跟著釋放，被取消的模型呼叫可能還沒真的停下來，和下一個呼叫短暫重疊（和改動前一樣）
    public func run<T>(_ op: () async -> T) async -> T {
        let t0 = Date()
        await acquireLocal()
        let fd = await acquireFile()
        state.withLock { waited += Date().timeIntervalSince(t0); calls += 1 }
        let start = Date().timeIntervalSince1970
        let result = await op()
        let end = Date().timeIntervalSince1970
        if let logPath { appendLog(logPath, String(format: "%d\t%.3f\t%.3f\n", getpid(), start, end)) }
        // 寫完日誌才釋放鎖：下一個呼叫的開始時間一定在這一行之後
        if fd >= 0 { flock(fd, LOCK_UN); close(fd) }
        releaseLocal()
        return result
    }

    private func acquireLocal() async {
        await withCheckedContinuation { (c: CheckedContinuation<Void, Never>) in
            state.lock()
            if busy { waiters.append(c); state.unlock() } else { busy = true; state.unlock(); c.resume() }
        }
    }

    // 輪到下一位時直接交棒（busy 維持 true），不讓後到的插隊
    private func releaseLocal() {
        state.lock()
        if waiters.isEmpty { busy = false; state.unlock() } else { let next = waiters.removeFirst(); state.unlock(); next.resume() }
    }

    // flock 會卡住執行緒，放到背景執行緒等，不卡 MainActor。開不了鎖檔就不跨程序排隊（只剩程序內）
    private func acquireFile() async -> Int32 {
        let path = lockPath
        return await withCheckedContinuation { c in
            DispatchQueue.global().async {
                let fd = open(path, O_CREAT | O_RDWR, 0o644)
                if fd >= 0 {
                    while flock(fd, LOCK_EX) != 0 && errno == EINTR {}
                }
                c.resume(returning: fd)
            }
        }
    }

    private func appendLog(_ path: String, _ line: String) {
        let fd = open(path, O_CREAT | O_WRONLY | O_APPEND, 0o644)
        guard fd >= 0 else { return }
        _ = line.withCString { write(fd, $0, strlen($0)) }
        close(fd)
    }
}

// 每個呼叫都經過排隊的模型（Bot 預設用它包 Apple 本機模型）
public struct QueuedModel: ChatModel {
    let base: ChatModel
    let queue: AIQueue
    public init(base: ChatModel, queue: AIQueue = .shared) { self.base = base; self.queue = queue }
    public func respond(_ request: ModelRequest) async -> ModelReply? {
        await queue.run { await base.respond(request) }
    }
}
