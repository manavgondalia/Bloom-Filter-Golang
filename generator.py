# generate training and testing txt files for the bloom filter from "words.txt"

import random


def generate_training_testing_files():
    with open("words.txt", "r") as f:
        words = f.readlines()
    words = [word.strip() for word in words]
    random.shuffle(words)
    training_words = words[:int(len(words) * 0.8)]
    testing_words = words[int(len(words) * 0.8):]
    with open("training.txt", "w") as f:
        for word in training_words:
            f.write(word + "\n")
    with open("testing.txt", "w") as f:
        for word in testing_words:
            f.write(word + "\n")
    print("Training and testing files generated successfully.")


if __name__ == "__main__":
    generate_training_testing_files()
