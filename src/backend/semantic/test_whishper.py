import os
import sys
import subprocess
import whisper
import torch


def extract_audio(filename):
    # Generate the audio file name based on the input file name
    audio_filename = filename.rsplit('.', 1)[0] + '_extracted.wav'

    # Check if the input file is already a .wav file
    if filename.endswith('.wav'):
        return filename

    # Check if the audio file already exists and delete it if it does
    if os.path.exists(audio_filename):
        os.remove(audio_filename)

    # Use FFmpeg to extract the audio
    subprocess.run(['ffmpeg', '-i', filename, '-q:a', '0', '-map', 'a', audio_filename], check=True)

    return audio_filename

def transcribe_audio(audio_filename):
    # Detect the available device
    device = "cpu"
    if torch.cuda.is_available():
        device = "cuda"
    elif torch.backends.mps.is_available() and torch.backends.mps.is_built():
        device = "mps"

    print(f"Using device: {device}")

    try:
        # Load the Whisper model on the detected device
        model = whisper.load_model("base", device=device).to(device)

        # Transcribe the audio file
        result = model.transcribe(audio_filename)
    except Exception as e:
        print(f"Encountered an error with device {device}: {e}")
        if device != "cpu":
            print("Falling back to CPU.")
            device = "cpu"
            model = whisper.load_model("base", device=device)
            result = model.transcribe(audio_filename)
        else:
            raise e

    return result["text"]

def extract_text_media(filename):
    # Extract audio from the media file
    audio_filename = extract_audio(filename)

    # Transcribe the extracted audio
    transcription = transcribe_audio(audio_filename)

    return transcription

if __name__ == "__main__":
    import pathlib

    try:
        filename = sys.argv[1]
    except IndexError:
        print(f"Usage:\npython {os.path.basename(__file__)} input_file")
        sys.exit()

    # Ensure the file exists
    if not os.path.exists(filename):
        print(f"File not found: {filename}")
        sys.exit()

    # Supported file extensions
    supported_extensions = {'mp4', 'mp3', 'mpeg', 'mpga', 'm4a', 'wav', 'webm', 'mov'}

    # Check if the file has a supported extension
    file_extension = pathlib.Path(filename).suffix[1:].lower()
    if file_extension not in supported_extensions:
        print(f"Unsupported file format: {file_extension}")
        print(f"Supported formats are: {', '.join(supported_extensions)}")
        sys.exit()

    # Extract and print the text from the media file
    print(extract_text_media(filename))
