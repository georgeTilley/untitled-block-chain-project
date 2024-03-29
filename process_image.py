import sys
import numpy as np
import tensorflow as tf
from tensorflow.keras.applications import InceptionV3
from tensorflow.keras.applications.inception_v3 import preprocess_input
from tensorflow.keras.preprocessing import image as keras_image
from tensorflow.keras.backend import sign

def preprocess_image(image_path):
    img = keras_image.load_img(image_path, target_size=(299, 299))
    img_array = keras_image.img_to_array(img)
    img_array = np.expand_dims(img_array, axis=0)
    img_array = preprocess_input(img_array)
    img_tensor = tf.convert_to_tensor(img_array)  # Convert to TensorFlow tensor
    return img_tensor

def preprocess_target(original_label):
    num_classes = 1000  # Assuming there are 1000 classes (adjust according to your problem)
    target = tf.one_hot(original_label, num_classes)
    target = tf.expand_dims(target, axis=0)  # Expand dims to match output shape
    return target

def adversarial_attack(image_path, epsilon=0.1):
    image = preprocess_image(image_path)
    pretrained_model = InceptionV3()

    # Obtain the original label (for demonstration purposes, you need to replace this with your actual label)
    original_label = 0  # For example, assume the original label is 0

    # Preprocess the target tensor
    target = preprocess_target(original_label)

    # Define the loss function
    loss_object = tf.keras.losses.CategoricalCrossentropy()

    with tf.GradientTape() as tape:
        tape.watch(image)
        prediction = pretrained_model(image)

        # Compute the loss
        loss = loss_object(target, prediction)

    # Calculate the gradient of the loss with respect to the input image
    gradient = tape.gradient(loss, image)
    
    # Get the sign of the gradient
    signed_grad = sign(gradient)

    # Create the adversarial image
    adversarial_image = image + epsilon * signed_grad

    # Clip the values to [0, 1]
    adversarial_image = tf.clip_by_value(adversarial_image, 0, 1)

    # Get the predicted label for the adversarial image
    adversarial_prediction = pretrained_model(adversarial_image)
    adversarial_label = tf.argmax(adversarial_prediction[0])

    return adversarial_image, original_label, adversarial_label

def save_adversarial_image(adversarial_image, output_path):
    img = keras_image.array_to_img(adversarial_image[0])
    img.save(output_path)
    print(output_path)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python adversarial_attack.py <image_path>")
        sys.exit(1)

    image_path = sys.argv[1]
    adversarial_image, original_label, adversarial_label = adversarial_attack(image_path)

    save_adversarial_image(adversarial_image, "adversarial_image.jpg")
