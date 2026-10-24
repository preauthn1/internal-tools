// meetnote-draft: instant draft transcript with macOS SpeechAnalyzer, shown in grey while recording
// until Breeze's accurate text replaces it. Built with swiftc from Command Line Tools (see build.rs).
//
// stdin:  raw 16 kHz mono float32 little-endian samples, from the start of the recording; EOF ends it.
// stdout: one JSON object per line:
//   {"kind":"ready"}
//   {"kind":"volatile","start":s,"end":s,"text":"…"}   still changing
//   {"kind":"final","start":s,"end":s,"text":"…"}      will not change
//   {"kind":"error","message":"…"}                      then exits with status 1
// Times are seconds of audio received.
import AVFoundation
import Foundation
import Speech

func emit(_ event: [String: Any]) {
    guard let data = try? JSONSerialization.data(withJSONObject: event),
          let line = String(data: data, encoding: .utf8) else { return }
    FileHandle.standardOutput.write((line + "\n").data(using: .utf8)!)
}

func fail(_ message: String) -> Never {
    emit(["kind": "error", "message": message])
    exit(1)
}

/// Reads stdin on its own thread and hands out analyzer-ready buffers.
func inputs(analyzerFormat: AVAudioFormat) -> AsyncStream<AnalyzerInput> {
    let source = AVAudioFormat(commonFormat: .pcmFormatFloat32, sampleRate: 16_000, channels: 1, interleaved: false)!
    let converter = analyzerFormat == source ? nil : AVAudioConverter(from: source, to: analyzerFormat)
    let (stream, continuation) = AsyncStream<AnalyzerInput>.makeStream()

    Thread {
        var carry = Data()
        while true {
            let chunk = FileHandle.standardInput.availableData
            if chunk.isEmpty { break }  // EOF: the recording stopped
            carry.append(chunk)
            let usable = carry.count - carry.count % 4
            guard usable > 0 else { continue }
            let frames = AVAudioFrameCount(usable / 4)
            guard let buffer = AVAudioPCMBuffer(pcmFormat: source, frameCapacity: frames) else { continue }
            buffer.frameLength = frames
            carry.prefix(usable).withUnsafeBytes { raw in
                buffer.floatChannelData![0].update(from: raw.bindMemory(to: Float.self).baseAddress!, count: Int(frames))
            }
            carry.removeFirst(usable)

            guard let converter else {
                continuation.yield(AnalyzerInput(buffer: buffer))
                continue
            }
            let capacity = AVAudioFrameCount(Double(frames) * analyzerFormat.sampleRate / 16_000) + 32
            guard let out = AVAudioPCMBuffer(pcmFormat: analyzerFormat, frameCapacity: capacity) else { continue }
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
            if out.frameLength > 0 { continuation.yield(AnalyzerInput(buffer: out)) }
        }
        continuation.finish()
    }.start()
    return stream
}

@main
struct Draft {
    static func main() async {
        let identifier = CommandLine.arguments.dropFirst().first ?? "zh_TW"
        guard SpeechTranscriber.isAvailable,
              let locale = await SpeechTranscriber.supportedLocale(equivalentTo: Locale(identifier: identifier))
        else { fail("這台 Mac 不支援 \(identifier) 的即時辨識") }

        let transcriber = SpeechTranscriber(
            locale: locale,
            transcriptionOptions: [],
            reportingOptions: [.volatileResults, .fastResults],
            attributeOptions: [.audioTimeRange]
        )
        do {
            if let request = try await AssetInventory.assetInstallationRequest(supporting: [transcriber]) {
                try await request.downloadAndInstall()
            }
        } catch {
            fail("無法準備即時辨識模型：\(error.localizedDescription)")
        }
        guard let format = await SpeechAnalyzer.bestAvailableAudioFormat(compatibleWith: [transcriber]) else {
            fail("無法取得即時辨識的音訊格式")
        }

        let analyzer = SpeechAnalyzer(modules: [transcriber])
        let results = Task {
            for try await result in transcriber.results {
                let text = String(result.text.characters).trimmingCharacters(in: .whitespaces)
                emit([
                    "kind": result.isFinal ? "final" : "volatile",
                    "start": result.range.start.seconds,
                    "end": result.range.end.seconds,
                    "text": text,
                ])
            }
        }
        do {
            try await analyzer.start(inputSequence: inputs(analyzerFormat: format))
            emit(["kind": "ready"])
            // Returns once stdin closes and every input has been analyzed.
            try await analyzer.finalizeAndFinishThroughEndOfInput()
            _ = try await results.value
        } catch {
            fail("即時辨識失敗：\(error.localizedDescription)")
        }
    }
}
