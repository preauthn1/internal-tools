// meetnote-tap: captures what the Mac is playing (the other side of a Google Meet, Slack or Zoom call)
// with a Core Audio process tap (macOS 14.2+). Built with swiftc from Command Line Tools (see build.rs).
//
// stdout: raw 16 kHz mono float32 little-endian samples, continuously (silence included) until stdin
//         closes or the process is terminated.
// stderr: one JSON object per line: {"kind":"ready","rate":48000} or {"kind":"error","message":"…"}.
//
// The first run makes macOS ask for "System Audio Recording" permission for the app. If it is denied,
// the tap keeps delivering silence rather than failing.
import AVFoundation
import CoreAudio
import Foundation

func report(_ event: [String: Any]) {
    guard let data = try? JSONSerialization.data(withJSONObject: event),
          let line = String(data: data, encoding: .utf8) else { return }
    FileHandle.standardError.write((line + "\n").data(using: .utf8)!)
}

func fail(_ message: String, _ status: OSStatus = 0) -> Never {
    report(["kind": "error", "message": status == 0 ? message : "\(message)（\(status)）"])
    exit(1)
}

func property<T>(_ object: AudioObjectID, _ selector: AudioObjectPropertySelector, _ value: inout T) -> OSStatus {
    var address = AudioObjectPropertyAddress(mSelector: selector, mScope: kAudioObjectPropertyScopeGlobal, mElement: kAudioObjectPropertyElementMain)
    var size = UInt32(MemoryLayout<T>.size)
    return AudioObjectGetPropertyData(object, &address, 0, nil, &size, &value)
}

func defaultOutputUID() -> String? {
    var device = AudioObjectID(kAudioObjectUnknown)
    guard property(AudioObjectID(kAudioObjectSystemObject), kAudioHardwarePropertyDefaultOutputDevice, &device) == noErr else { return nil }
    var address = AudioObjectPropertyAddress(mSelector: kAudioDevicePropertyDeviceUID, mScope: kAudioObjectPropertyScopeGlobal, mElement: kAudioObjectPropertyElementMain)
    var uid: Unmanaged<CFString>?
    var size = UInt32(MemoryLayout<Unmanaged<CFString>?>.size)
    guard AudioObjectGetPropertyData(device, &address, 0, nil, &size, &uid) == noErr, let uid else { return nil }
    return uid.takeRetainedValue() as String
}

// 1. A private mono tap on everything the system plays.
let description = CATapDescription(monoGlobalTapButExcludeProcesses: [])
description.uuid = UUID()
description.isPrivate = true
description.muteBehavior = .unmuted
var tap = AudioObjectID(kAudioObjectUnknown)
var status = AudioHardwareCreateProcessTap(description, &tap)
guard status == noErr else { fail("無法建立系統聲音擷取", status) }

var tapFormat = AudioStreamBasicDescription()
status = property(tap, kAudioTapPropertyFormat, &tapFormat)
guard status == noErr, let sourceFormat = AVAudioFormat(streamDescription: &tapFormat) else { fail("無法讀取系統聲音格式", status) }

// 2. Taps are read through a private aggregate device, clocked by the current output device.
guard let outputUID = defaultOutputUID() else { fail("找不到聲音輸出裝置") }
let aggregateDescription: [String: Any] = [
    kAudioAggregateDeviceNameKey: "meetnote 系統聲音",
    kAudioAggregateDeviceUIDKey: UUID().uuidString,
    kAudioAggregateDeviceMainSubDeviceKey: outputUID,
    kAudioAggregateDeviceIsPrivateKey: true,
    kAudioAggregateDeviceIsStackedKey: false,
    kAudioAggregateDeviceTapAutoStartKey: true,
    kAudioAggregateDeviceSubDeviceListKey: [[kAudioSubDeviceUIDKey: outputUID]],
    kAudioAggregateDeviceTapListKey: [[kAudioSubTapDriftCompensationKey: true, kAudioSubTapUIDKey: description.uuid.uuidString]],
]
var aggregate = AudioObjectID(kAudioObjectUnknown)
status = AudioHardwareCreateAggregateDevice(aggregateDescription as CFDictionary, &aggregate)
guard status == noErr else { fail("無法建立系統聲音裝置", status) }

@Sendable func cleanUp() {
    AudioHardwareDestroyAggregateDevice(aggregate)
    AudioHardwareDestroyProcessTap(tap)
}

// 3. Convert every buffer to 16 kHz mono float32 and write it out.
let target = AVAudioFormat(commonFormat: .pcmFormatFloat32, sampleRate: 16_000, channels: 1, interleaved: false)!
guard let converter = AVAudioConverter(from: sourceFormat, to: target) else {
    cleanUp()
    fail("無法轉換系統聲音格式")
}
let stdout = FileHandle.standardOutput
let queue = DispatchQueue(label: "meetnote.tap")
var procID: AudioDeviceIOProcID?
status = AudioDeviceCreateIOProcIDWithBlock(&procID, aggregate, queue) { _, inputData, _, _, _ in
    guard let input = AVAudioPCMBuffer(pcmFormat: sourceFormat, bufferListNoCopy: inputData, deallocator: nil),
          input.frameLength > 0 else { return }
    let capacity = AVAudioFrameCount(Double(input.frameLength) * 16_000 / sourceFormat.sampleRate) + 32
    guard let out = AVAudioPCMBuffer(pcmFormat: target, frameCapacity: capacity) else { return }
    var supplied = false
    _ = converter.convert(to: out, error: nil) { _, inputStatus in
        if supplied {
            inputStatus.pointee = .noDataNow
            return nil
        }
        supplied = true
        inputStatus.pointee = .haveData
        return input
    }
    guard out.frameLength > 0 else { return }
    let bytes = Data(bytes: out.floatChannelData![0], count: Int(out.frameLength) * 4)
    do {
        try stdout.write(contentsOf: bytes)
    } catch {
        // The app went away; stop quietly.
        exit(0)
    }
}
guard status == noErr, let procID else {
    cleanUp()
    fail("無法開始擷取系統聲音", status)
}
status = AudioDeviceStart(aggregate, procID)
guard status == noErr else {
    AudioDeviceDestroyIOProcID(aggregate, procID)
    cleanUp()
    fail("無法開始擷取系統聲音", status)
}
report(["kind": "ready", "rate": sourceFormat.sampleRate, "channels": sourceFormat.channelCount])

// 4. Run until the app closes our stdin or terminates us; always remove the tap on the way out.
func stop() -> Never {
    AudioDeviceStop(aggregate, procID)
    AudioDeviceDestroyIOProcID(aggregate, procID)
    cleanUp()
    exit(0)
}
signal(SIGTERM, SIG_IGN)
let term = DispatchSource.makeSignalSource(signal: SIGTERM, queue: .main)
term.setEventHandler { stop() }
term.resume()
Thread {
    while !FileHandle.standardInput.availableData.isEmpty {}
    DispatchQueue.main.async { stop() }
}.start()
dispatchMain()
