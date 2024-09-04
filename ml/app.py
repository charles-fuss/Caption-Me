
from transformers import VisionEncoderDecoderModel, ViTImageProcessor, AutoTokenizer
import torch
import ast
from PIL import Image
from flask import Flask, request

model = VisionEncoderDecoderModel.from_pretrained("nlpconnect/vit-gpt2-image-captioning")
feature_extractor = ViTImageProcessor.from_pretrained("nlpconnect/vit-gpt2-image-captioning")
tokenizer = AutoTokenizer.from_pretrained("nlpconnect/vit-gpt2-image-captioning")

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
model.to(device)


app = Flask(__name__)

@app.route('/caption_image', methods=['POST'])
def caption_image():
    breakpoint()
    data = ast.literal_eval(request.data.decode('utf-8'))
    if 'image' in data.keys():
        image_binary = data['image']
    else:
        return {"error": "no image key, invalid keys sent"}
    max_length = 16
    num_beams = 4
    gen_kwargs = {"max_length": max_length, "num_beams": num_beams}

    try:
        i_image = Image.open(image_binary)
    except:
        return {"error": "not valid image binary"}
    if i_image.mode != "RGB":
        i_image = i_image.convert(mode="RGB")
    else:
        return {'error', 'not RGB'}
    pixel_values = feature_extractor(images=i_image, return_tensors="pt").pixel_values
    pixel_values = pixel_values.to(device)

    output_ids = model.generate(pixel_values, **gen_kwargs)

    preds = tokenizer.batch_decode(output_ids, skip_special_tokens=True)
    preds = [pred.strip() for pred in preds]
    return preds


