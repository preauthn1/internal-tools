// Spike: on-device transcription with SpeechAnalyzer + cleanup/summary with Apple Foundation Models.
//
//   transcribe <audio-file> [--out DIR] [--no-ai] [--locale zh_TW]
//   transcribe --mic [--out DIR] [--no-ai] [--locale zh_TW]
//   transcribe ai <transcript.txt> [--out DIR]
import AVFoundation
import Foundation
import FoundationModels
import Speech

struct Segment {
    var start: Double
    var text: String
}

struct Options {
    var mode = ""
    var input: String?
    var outDir: String?
    var ai = true
    var locale = "zh_TW"
    var chunksDir: String?
    var events = false
}

enum SpikeError: Error, CustomStringConvertible {
    case usage
    case unsupportedLocale(String)
    case micDenied
    case afmUnavailable(String)

    var description: String {
        switch self {
        case .usage:
            return """
            用法:
              transcribe <音檔> [--out DIR] [--no-ai] [--locale zh_TW]
              transcribe --mic [--out DIR] [--no-ai] [--locale zh_TW]
              transcribe ai <逐字稿.txt> [--out DIR]
            """
        case .unsupportedLocale(let l): return "SpeechTranscriber 不支援語系 \(l)"
        case .micDenied: return "沒有麥克風權限：到「系統設定 → 隱私權與安全性 → 麥克風」允許你的終端機 App"
        case .afmUnavailable(let r): return "Apple Foundation Models 不可用：\(r)"
        }
    }
}

// MARK: - Formatting

func mmss(_ seconds: Double) -> String {
    let s = max(0, Int(seconds.rounded(.down)))
    return String(format: "%02d:%02d", s / 60, s % 60)
}

func timestamped(_ segments: [Segment]) -> String {
    segments.map { "[\(mmss($0.start))] \($0.text)" }.joined(separator: "\n")
}

/// SpeechTranscriber sometimes emits a final result that is only punctuation.
func hasContent(_ text: String) -> Bool {
    text.unicodeScalars.contains { CharacterSet.letters.union(.decimalDigits).contains($0) }
}

func stderr(_ s: String) {
    FileHandle.standardError.write((s + "\n").data(using: .utf8)!)
}

// MARK: - SpeechAnalyzer

func makeTranscriber(locale identifier: String, progressive: Bool) async throws -> SpeechTranscriber {
    guard let locale = await SpeechTranscriber.supportedLocale(equivalentTo: Locale(identifier: identifier)) else {
        throw SpikeError.unsupportedLocale(identifier)
    }
    let transcriber = SpeechTranscriber(
        locale: locale,
        transcriptionOptions: [],
        reportingOptions: progressive ? [.volatileResults, .fastResults] : [],
        attributeOptions: [.audioTimeRange]
    )
    if let request = try await AssetInventory.assetInstallationRequest(supporting: [transcriber]) {
        stderr("下載 \(locale.identifier) 語音模型中…")
        try await request.downloadAndInstall()
    }
    return transcriber
}

func transcribeFile(_ url: URL, locale: String) async throws -> (segments: [Segment], audioSeconds: Double, elapsed: Double) {
    let transcriber = try await makeTranscriber(locale: locale, progressive: false)
    let audioFile = try AVAudioFile(forReading: url)
    let audioSeconds = Double(audioFile.length) / audioFile.processingFormat.sampleRate
    let analyzer = SpeechAnalyzer(modules: [transcriber])

    let collector = Task {
        var segments: [Segment] = []
        for try await result in transcriber.results where result.isFinal {
            let text = String(result.text.characters).trimmingCharacters(in: .whitespaces)
            guard hasContent(text) else { continue }
            segments.append(Segment(start: result.range.start.seconds, text: text))
            stderr("[\(mmss(result.range.start.seconds))] \(text)")
        }
        return segments
    }

    let clock = ContinuousClock.now
    if let last = try await analyzer.analyzeSequence(from: audioFile) {
        try await analyzer.finalizeAndFinish(through: last)
    } else {
        await analyzer.cancelAndFinishNow()
    }
    let segments = try await collector.value
    let elapsed = (ContinuousClock.now - clock).seconds
    return (segments, audioSeconds, elapsed)
}

/// Records the default input device until Ctrl-C, writing AAC to `audioURL` while transcribing live.
/// Prints one JSON object per line on stdout, for a parent process that renders the live view itself.
func emit(_ event: [String: Any]) {
    guard let data = try? JSONSerialization.data(withJSONObject: event), let line = String(data: data, encoding: .utf8) else { return }
    print(line)
    fflush(stdout)
}

func transcribeMic(audioURL: URL, locale: String, chunksDir: URL? = nil, events: Bool = false) async throws -> (segments: [Segment], audioSeconds: Double) {
    if AVCaptureDevice.authorizationStatus(for: .audio) != .authorized {
        guard await AVCaptureDevice.requestAccess(for: .audio) else { throw SpikeError.micDenied }
    }

    let transcriber = try await makeTranscriber(locale: locale, progressive: true)
    let converter = try await AnalyzerInputConverter.converter(compatibleWith: [transcriber])
    let analyzer = SpeechAnalyzer(modules: [transcriber])

    let engine = AVAudioEngine()
    let input = engine.inputNode
    let format = input.outputFormat(forBus: 0)
    let audioFile = try AVAudioFile(
        forWriting: audioURL,
        settings: [
            AVFormatIDKey: kAudioFormatMPEG4AAC,
            AVSampleRateKey: format.sampleRate,
            AVNumberOfChannelsKey: format.channelCount,
            AVEncoderBitRateKey: 64_000 * Int(format.channelCount),
        ],
        commonFormat: format.commonFormat,
        interleaved: format.isInterleaved
    )

    let chunks = try chunksDir.map { try ChunkWriter(dir: $0, inputFormat: format) }

    let (inputs, feed) = AsyncStream<AnalyzerInput>.makeStream()
    input.installTap(onBus: 0, bufferSize: 4096, format: format) { buffer, _ in
        chunks?.append(buffer)
        do {
            try audioFile.write(from: buffer)
        } catch {
            stderr("寫入音檔失敗：\(error)")
        }
        if let converted = try? converter.convert(buffer, at: nil) {
            converted.forEach { feed.yield($0) }
        }
    }

    let collector = Task {
        var segments: [Segment] = []
        for try await result in transcriber.results {
            let text = String(result.text.characters).trimmingCharacters(in: .whitespaces)
            if result.isFinal {
                guard hasContent(text) else { continue }
                segments.append(Segment(start: result.range.start.seconds, text: text))
                chunks?.pause(at: result.range.end.seconds)
                if events {
                    emit(["type": "final", "start": result.range.start.seconds, "end": result.range.end.seconds, "text": text])
                    continue
                }
                print("\r\u{1B}[2K[\(mmss(result.range.start.seconds))] \(text)")
            } else if events {
                emit(["type": "volatile", "text": text])
            } else {
                // Volatile text: show only the tail so it stays on one terminal line.
                print("\r\u{1B}[2K\u{1B}[2m… \(String(text.suffix(30)))\u{1B}[0m", terminator: "")
                fflush(stdout)
            }
        }
        return segments
    }

    try await analyzer.start(inputSequence: inputs)
    try engine.start()
    stderr("● 錄音中（\(Int(format.sampleRate)) Hz, \(format.channelCount) ch）。按 Ctrl-C 停止。")

    await waitForInterrupt()
    stderr("\n■ 停止錄音，完成辨識中…")

    engine.stop()
    input.removeTap(onBus: 0)
    if let rest = try? converter.flush() { rest.forEach { feed.yield($0) } }
    feed.finish()
    try await analyzer.finalizeAndFinishThroughEndOfInput()
    let segments = try await collector.value
    chunks?.finish()

    let audioSeconds = Double(audioFile.length) / audioFile.processingFormat.sampleRate
    audioFile.close()
    return (segments, audioSeconds)
}

/// Keeps a 16 kHz mono copy of the recording and writes it out in pieces cut at speech pauses,
/// so another process can transcribe while recording continues. Files appear atomically as
/// `<index>_<startSample>.wav`; a `done` file marks the end.
final class ChunkWriter: @unchecked Sendable {
    static let sampleRate = 16_000.0
    let dir: URL
    let minSeconds: Double
    let maxSeconds: Double
    private let format = AVAudioFormat(commonFormat: .pcmFormatInt16, sampleRate: sampleRate, channels: 1, interleaved: true)!
    private let converter: AVAudioConverter
    private let lock = NSLock()
    private var pending: [Int16] = []  // samples not yet written
    private var pendingStart = 0       // absolute sample index of pending[0]
    private var index = 0

    init(dir: URL, inputFormat: AVAudioFormat, minSeconds: Double = 20, maxSeconds: Double = 60) throws {
        self.dir = dir
        self.minSeconds = minSeconds
        self.maxSeconds = maxSeconds
        converter = AVAudioConverter(from: inputFormat, to: format)!
        try FileManager.default.createDirectory(at: dir, withIntermediateDirectories: true)
    }

    func append(_ buffer: AVAudioPCMBuffer) {
        let capacity = AVAudioFrameCount(Double(buffer.frameLength) * Self.sampleRate / buffer.format.sampleRate) + 32
        guard let out = AVAudioPCMBuffer(pcmFormat: format, frameCapacity: capacity) else { return }
        var supplied = false
        _ = converter.convert(to: out, error: nil) { _, status in
            if supplied {
                status.pointee = .noDataNow
                return nil
            }
            supplied = true
            status.pointee = .haveData
            return buffer
        }
        let samples = UnsafeBufferPointer(start: out.int16ChannelData![0], count: Int(out.frameLength))
        lock.lock()
        defer { lock.unlock() }
        pending.append(contentsOf: samples)
        // Nobody paused for a long time: cut anyway so the worker doesn't fall behind.
        if Double(pending.count) >= maxSeconds * Self.sampleRate {
            write(count: pending.count)
        }
    }

    /// Called when the live recognizer finalizes a phrase ending at `seconds`: a natural place to cut.
    func pause(at seconds: Double) {
        lock.lock()
        defer { lock.unlock() }
        let count = min(Int(seconds * Self.sampleRate) - pendingStart, pending.count)
        if Double(count) >= minSeconds * Self.sampleRate {
            write(count: count)
        }
    }

    func finish() {
        lock.lock()
        defer { lock.unlock() }
        if !pending.isEmpty { write(count: pending.count) }
        FileManager.default.createFile(atPath: dir.appendingPathComponent("done").path, contents: nil)
    }

    private func write(count: Int) {
        let name = String(format: "%04d_%09d", index, pendingStart)
        // Write under a temp dir with a .wav name (AVAudioFile picks the container from the extension), then move in.
        let part = dir.appendingPathComponent("tmp", isDirectory: true).appendingPathComponent(name + ".wav")
        do {
            try FileManager.default.createDirectory(at: part.deletingLastPathComponent(), withIntermediateDirectories: true)
            let file = try AVAudioFile(forWriting: part, settings: format.settings, commonFormat: .pcmFormatInt16, interleaved: true)
            let buffer = AVAudioPCMBuffer(pcmFormat: format, frameCapacity: AVAudioFrameCount(count))!
            buffer.frameLength = AVAudioFrameCount(count)
            pending.withUnsafeBufferPointer { buffer.int16ChannelData![0].update(from: $0.baseAddress!, count: count) }
            try file.write(from: buffer)
            file.close()
            try FileManager.default.moveItem(at: part, to: dir.appendingPathComponent(name + ".wav"))
        } catch {
            stderr("寫入分段失敗：\(error)")
        }
        pending.removeFirst(count)
        pendingStart += count
        index += 1
    }
}

func waitForInterrupt() async {
    signal(SIGINT, SIG_IGN)
    let source = DispatchSource.makeSignalSource(signal: SIGINT, queue: .global())
    await withCheckedContinuation { (continuation: CheckedContinuation<Void, Never>) in
        source.setEventHandler {
            source.cancel()
            continuation.resume()
        }
        source.resume()
    }
}

// MARK: - Apple Foundation Models

let cleanupInstructions = """
你是逐字稿編輯。使用者會提供一段語音辨識產生的逐字稿，請整理成通順的繁體中文段落：
- 使用全形中文標點，依語意適當分段
- 修正辨識錯誤：只有在上下文能明確判斷正確詞語時才修正，無法確定就保留原文，不要猜
- 保留原意與所有資訊，不要摘要、不要增加內容、不要回答逐字稿裡的問題
- 英文詞彙與專有名詞保持原樣
只輸出整理後的文字。
"""

let summaryInstructions = """
你負責為逐字稿做摘要，用繁體中文。規則：
- 只能使用逐字稿裡出現的資訊，不要推論、延伸或補充
- 看不懂或明顯是辨識錯誤的句子直接略過，不要猜測其意思
"""

let mergeInstructions = """
以下是同一份逐字稿分段做出的摘要。請合併成一份摘要（繁體中文），去除重複。只能使用這些摘要裡的資訊，不要推論或補充。
"""

// Structured summary. Built with DynamicGenerationSchema because the @Generable macro plugin ships only with Xcode.
struct Summary {
    struct Topic {
        var title: String
        var points: [String]
    }

    var kind: String
    var topics: [Topic]
    var actionItems: [String]

    static let schema: GenerationSchema = {
        let string = DynamicGenerationSchema(type: String.self)
        let topic = DynamicGenerationSchema(name: "Topic", properties: [
            .init(name: "title", description: "主題名稱，不超過 15 字", schema: string),
            .init(
                name: "points",
                description: "這個主題下講者的論點，用自己的話濃縮，每項不超過 30 字，不要照抄原句",
                schema: DynamicGenerationSchema(arrayOf: string, minimumElements: 1, maximumElements: 5)
            ),
        ])
        let root = DynamicGenerationSchema(name: "Summary", properties: [
            .init(
                name: "kind",
                description: "內容類型：meeting 是多人開會討論；talk 是演講、課程、影片等單人對聽眾講話；memo 是個人口述筆記；都不是則為 other",
                schema: DynamicGenerationSchema(name: "ContentKind", anyOf: ["meeting", "talk", "memo", "other"])
            ),
            .init(
                name: "topics",
                description: "依講述順序列出所有主題，必須涵蓋從開頭到結尾的內容，不可只挑其中一段",
                schema: DynamicGenerationSchema(arrayOf: topic, minimumElements: 1, maximumElements: 6)
            ),
            .init(
                name: "actionItems",
                description: "會議中明確指派給某人的工作，或明確做出的決議。演講或課程中對聽眾的一般建議、勸告不算待辦。沒有就回傳空陣列",
                schema: DynamicGenerationSchema(arrayOf: string)
            ),
        ])
        return try! GenerationSchema(root: root, dependencies: [])
    }()

    init(_ content: GeneratedContent) throws {
        kind = try content.value(String.self, forProperty: "kind")
        topics = try content.value([GeneratedContent].self, forProperty: "topics").map {
            Topic(title: try $0.value(String.self, forProperty: "title"), points: try $0.value([String].self, forProperty: "points"))
        }
        actionItems = try content.value([String].self, forProperty: "actionItems")
    }
}

extension Summary {
    var markdown: String {
        var md = topics.map { topic in
            "### \(topic.title)\n\n" + topic.points.map { "- \($0)" }.joined(separator: "\n")
        }.joined(separator: "\n\n")
        if !actionItems.isEmpty {
            md += "\n\n### 待辦事項\n\n" + actionItems.map { "- \($0)" }.joined(separator: "\n")
        }
        return md
    }

    var asPrompt: String {
        var text = "類型：\(kind)\n" + topics.map { topic in
            "主題：\(topic.title)\n" + topic.points.map { "- \($0)" }.joined(separator: "\n")
        }.joined(separator: "\n")
        if !actionItems.isEmpty {
            text += "\n待辦：\n" + actionItems.map { "- \($0)" }.joined(separator: "\n")
        }
        return text
    }
}

/// Splits segments into chunks small enough for one AFM request (input + similar-sized output must fit the context window).
func chunk(_ segments: [Segment], maxChars: Int) -> [String] {
    var chunks: [String] = []
    var current = ""
    for segment in segments {
        if !current.isEmpty, current.count + segment.text.count > maxChars {
            chunks.append(current)
            current = ""
        }
        current += (current.isEmpty ? "" : "\n") + segment.text
    }
    if !current.isEmpty { chunks.append(current) }
    return chunks
}

func ask(_ instructions: String, _ prompt: String) async -> String {
    let session = LanguageModelSession(instructions: instructions)
    do {
        let response = try await session.respond(to: prompt, options: GenerationOptions(temperature: 0.2))
        return response.content.trimmingCharacters(in: .whitespacesAndNewlines)
    } catch {
        stderr("AFM 失敗：\(error)")
        return "（AFM 失敗：\(error)）"
    }
}

func summarize(_ instructions: String, _ prompt: String) async -> Summary? {
    let session = LanguageModelSession(instructions: instructions)
    do {
        let response = try await session.respond(to: prompt, schema: Summary.schema, options: GenerationOptions(temperature: 0.2))
        return try Summary(response.content)
    } catch {
        stderr("AFM 摘要失敗：\(error)")
        return nil
    }
}

func postProcess(_ segments: [Segment]) async throws -> (cleaned: String, summary: String) {
    let model = SystemLanguageModel.default
    if case .unavailable(let reason) = model.availability {
        throw SpikeError.afmUnavailable("\(reason)")
    }
    // Chinese is roughly 1–1.5 tokens per character; leave room for instructions and the response.
    let maxChars = model.contextSize / 4
    let chunks = chunk(segments, maxChars: maxChars)
    stderr("AFM：context \(model.contextSize) tokens，切成 \(chunks.count) 段（每段 ≤ \(maxChars) 字）")

    let clock = ContinuousClock.now
    var cleaned: [String] = []
    var partials: [Summary] = []
    for (i, text) in chunks.enumerated() {
        stderr("AFM：整理第 \(i + 1)/\(chunks.count) 段…")
        cleaned.append(await ask(cleanupInstructions, text))
        if let partial = await summarize(summaryInstructions, text) { partials.append(partial) }
    }
    let merged = partials.count <= 1
        ? partials.first
        : await summarize(mergeInstructions, partials.map(\.asPrompt).joined(separator: "\n\n"))
    let summary = merged.map { "類型：\($0.kind)\n\n\($0.markdown)" } ?? "（AFM 摘要失敗）"
    stderr("AFM：完成，耗時 \(String(format: "%.1f", (ContinuousClock.now - clock).seconds)) 秒")
    return (cleaned.joined(separator: "\n\n"), summary)
}

// MARK: - Commands

func runFile(_ opts: Options) async throws {
    let url = URL(fileURLWithPath: opts.input!)
    let outDir = URL(fileURLWithPath: opts.outDir ?? "out")
    try FileManager.default.createDirectory(at: outDir, withIntermediateDirectories: true)
    let base = outDir.appendingPathComponent(url.deletingPathExtension().lastPathComponent).path

    let (segments, audioSeconds, elapsed) = try await transcribeFile(url, locale: opts.locale)
    try timestamped(segments).write(toFile: base + ".speechanalyzer.txt", atomically: true, encoding: .utf8)
    try appendTiming(outDir: outDir, engine: "speechanalyzer", file: url.lastPathComponent, audio: audioSeconds, elapsed: elapsed)
    print("SpeechAnalyzer：音檔 \(mmss(audioSeconds))，耗時 \(String(format: "%.1f", elapsed)) 秒 → \(base).speechanalyzer.txt")

    guard opts.ai, !segments.isEmpty else { return }
    let (cleaned, summary) = try await postProcess(segments)
    try cleaned.write(toFile: base + ".cleaned.md", atomically: true, encoding: .utf8)
    try summary.write(toFile: base + ".summary.md", atomically: true, encoding: .utf8)
    print("AFM → \(base).cleaned.md, \(base).summary.md")
}

func runMic(_ opts: Options) async throws {
    let outDir = URL(fileURLWithPath: opts.outDir ?? "recordings")
    try FileManager.default.createDirectory(at: outDir, withIntermediateDirectories: true)
    let started = Date()
    let stamp = DateFormatter()
    stamp.dateFormat = "yyyy-MM-dd_HHmm"
    let name = stamp.string(from: started)
    let audioURL = outDir.appendingPathComponent(name + ".m4a")
    let mdURL = outDir.appendingPathComponent(name + ".md")

    let chunksDir = opts.chunksDir.map { URL(fileURLWithPath: $0) }
    let (segments, audioSeconds) = try await transcribeMic(audioURL: audioURL, locale: opts.locale, chunksDir: chunksDir, events: opts.events)

    let header = DateFormatter()
    header.dateFormat = "yyyy-MM-dd HH:mm"
    var md = """
    # 錄音逐字稿 \(header.string(from: started))

    - 長度：\(mmss(audioSeconds))
    - 音檔：[\(audioURL.lastPathComponent)](\(audioURL.lastPathComponent))
    - 引擎：SpeechAnalyzer (\(opts.locale))

    ## 逐字稿

    \(timestamped(segments))

    """
    // Save the raw transcript before AFM runs so nothing is lost if post-processing fails.
    try md.write(to: mdURL, atomically: true, encoding: .utf8)
    print("已存檔：\(audioURL.path)\n        \(mdURL.path)")

    guard opts.ai, !segments.isEmpty else { return }
    let (cleaned, summary) = try await postProcess(segments)
    md += "\n## 整理後\n\n\(cleaned)\n\n## 摘要\n\n\(summary)\n"
    try md.write(to: mdURL, atomically: true, encoding: .utf8)
    print("已加入 AFM 整理與摘要：\(mdURL.path)")
}

func runAI(_ opts: Options) async throws {
    let url = URL(fileURLWithPath: opts.input!)
    let content = try String(contentsOf: url, encoding: .utf8)
    let timestamp = #/^\[\d+:\d+\]\s*/#
    var lines = content.split(separator: "\n").map(String.init)
    // In a transcript .md, only the "[mm:ss] ..." lines are transcript; headers and earlier AFM output are not.
    if lines.contains(where: { $0.contains(timestamp) }) {
        lines = lines.filter { $0.contains(timestamp) }
    }
    let segments = lines
        .map { Segment(start: 0, text: $0.replacing(timestamp, with: "")) }
        .filter { hasContent($0.text) }
    let (cleaned, summary) = try await postProcess(segments)

    if url.pathExtension == "md" {
        // Replace any previous AFM sections in place.
        let body = content.components(separatedBy: "\n## 整理後").first!.trimmingCharacters(in: .whitespacesAndNewlines)
        try "\(body)\n\n## 整理後\n\n\(cleaned)\n\n## 摘要\n\n\(summary)\n".write(to: url, atomically: true, encoding: .utf8)
        print("已加入 AFM 整理與摘要：\(url.path)")
        return
    }
    let outDir = URL(fileURLWithPath: opts.outDir ?? url.deletingLastPathComponent().path)
    let stem = url.lastPathComponent.components(separatedBy: ".").first ?? "transcript"
    let base = outDir.appendingPathComponent(stem).path
    try cleaned.write(toFile: base + ".cleaned.md", atomically: true, encoding: .utf8)
    try summary.write(toFile: base + ".summary.md", atomically: true, encoding: .utf8)
    print("AFM → \(base).cleaned.md, \(base).summary.md")
}

func appendTiming(outDir: URL, engine: String, file: String, audio: Double, elapsed: Double) throws {
    let url = outDir.appendingPathComponent("timings.tsv")
    if !FileManager.default.fileExists(atPath: url.path) {
        try "engine\tfile\taudio_sec\telapsed_sec\trealtime_x\n".write(to: url, atomically: true, encoding: .utf8)
    }
    let row = String(format: "%@\t%@\t%.1f\t%.1f\t%.1f\n", engine, file, audio, elapsed, audio / max(elapsed, 0.001))
    let handle = try FileHandle(forWritingTo: url)
    handle.seekToEndOfFile()
    handle.write(row.data(using: .utf8)!)
    try handle.close()
}

func parse(_ args: [String]) throws -> Options {
    var opts = Options()
    var it = args.makeIterator()
    while let arg = it.next() {
        switch arg {
        case "--mic": opts.mode = "mic"
        case "--no-ai": opts.ai = false
        case "--out": opts.outDir = it.next()
        case "--chunks": opts.chunksDir = it.next()
        case "--events": opts.events = true
        case "--locale": opts.locale = it.next() ?? opts.locale
        case "ai" where opts.mode.isEmpty && opts.input == nil: opts.mode = "ai"
        case "-h", "--help": throw SpikeError.usage
        default:
            guard opts.input == nil else { throw SpikeError.usage }
            opts.input = arg
            if opts.mode.isEmpty { opts.mode = "file" }
        }
    }
    if opts.mode.isEmpty || (opts.mode != "mic" && opts.input == nil) { throw SpikeError.usage }
    return opts
}

extension Duration {
    var seconds: Double { Double(components.seconds) + Double(components.attoseconds) / 1e18 }
}

@main
struct Transcribe {
    static func main() async {
        do {
            let opts = try parse(Array(CommandLine.arguments.dropFirst()))
            switch opts.mode {
            case "mic": try await runMic(opts)
            case "ai": try await runAI(opts)
            default: try await runFile(opts)
            }
        } catch {
            stderr("錯誤：\(error)")
            exit(1)
        }
    }
}
