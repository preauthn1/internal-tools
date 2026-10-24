import Foundation

// MARK: - 伺服器日誌：一行一筆，輸出目的地可以換（測試會擷取下來檢查）
// 寫出去之前把 token、secret 換成 ***（最後一道防線；呼叫的地方本來就不該把它們放進日誌）

public final class Log: @unchecked Sendable {
    private let lock = NSLock()
    private let secrets: [String]
    private let write: @Sendable (String) -> Void

    public init(secrets: [String] = [], write: @escaping @Sendable (String) -> Void = { print($0); fflush(stdout) }) {
        self.secrets = secrets.filter { !$0.isEmpty }
        self.write = write
    }

    public func info(_ s: String) { emit("ℹ️ " + s) }
    public func warn(_ s: String) { emit("⚠️ " + s) }
    public func error(_ s: String) { emit("❌ " + s) }

    public func redact(_ s: String) -> String {
        secrets.reduce(s) { $0.replacingOccurrences(of: $1, with: "***") }
    }

    private func emit(_ s: String) {
        let line = "\(ISO8601DateFormatter().string(from: Date())) \(redact(s))"
        lock.withLock { write(line) }
    }
}
