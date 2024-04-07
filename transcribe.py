from pydub import AudioSegment
import time
from vosk import Model, KaldiRecognizer
import wave
import json
import wave 
import argparse

args = argparse.ArgumentParser()
args.add_argument("--input_audio", type=str)
args.add_argument("--vid_id", type=str)

input_audio_file_path = args.parse_args().input_audio
vid_id = args.parse_args().vid_id
# Load audio file
audio = AudioSegment.from_wav(input_audio_file_path)

# Convert to mono
audio = audio.set_channels(1)

# Set frame rate
audio = audio.set_frame_rate(16000)

# Save new audio
audio.export(f"cache/output_audio_mono-{vid_id}.wav", format="wav")



model = Model("vosk-model-small-en-us-0.15")
rec = KaldiRecognizer(model, 16000)

# Open WAV file
wf = wave.open(f"cache/output_audio_mono-{vid_id}.wav", "rb")

# List to hold all text segments
transcribed_text_list = []

start_time = time.time()
while True:
    if time.time() - start_time > 1 * 60:
        break
    data = wf.readframes(4000)
    if len(data) == 0:
        break
    if rec.AcceptWaveform(data):
        result = json.loads(rec.Result())
        print(f"result: {result}, {result['text']}")
        transcribed_text_list.append(result['text'])

# Handle last part of audio
final_result = json.loads(rec.FinalResult())
transcribed_text_list.append(final_result['text'])

# Concatenate all text segments
complete_text = ' '.join(transcribed_text_list)

# Write the complete transcribed text to a text file
with open(f"transcripts/transcribed_text-{vid_id}.txt", "w") as f:
    f.write(complete_text)

print("Transcription complete. Output written to transcripts/transcribed_text-{vid_id}.txt")