<?php

class Utils
{
    private const ID_LENGTH = 32;

    static public function genUuid(): string
    {
        return sprintf('%04x%04x%04x%04x%04x%04x%04x%04x',
            mt_rand(0, 0xffff), mt_rand(0, 0xffff),
            mt_rand(0, 0xffff),
            mt_rand(0, 0x0fff) | 0x4000,
            mt_rand(0, 0x3fff) | 0x8000,
            mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff)
        );
    }

    static public function isUuid(string $id): bool
    {
        return preg_match('/^[a-z0-9]{' . self::ID_LENGTH . '}$/', $id) === 1;
    }

    static public function isValidString(string $str): bool
    {
        return preg_match('/^[0-9a-zA-Zа-яА-Я\s\-\'\.\,]{1,64}$/', $str) === 1;

    }

    static public function randomName(): string
    {
        $adjectives = ['Amiable', 'Blissful', 'Cheerful', 'Delightful', 'Enchanting', 'Friendly', 'Gracious', 'Harmonious', 'Invigorating', 'Jovial', 'Kindhearted', 'Lively', 'Magnificent', 'Nurturing', 'Optimistic', 'Playful', 'Quaint', 'Radiant', 'Serene', 'Tranquil', 'Uplifting', 'Vibrant', 'Wholesome', 'Affectionate', 'Beautiful', 'Charming', 'Dreamy', 'Elegant', 'Festive', 'Gentle', 'Heart warming', 'Inspiring', 'Jubilant', 'Kind', 'Lovely', 'Majestic', 'Noble', 'Outstanding', 'Pleasurable', 'Radiant', 'Splendid', 'Tender', 'Unforgettable', 'Virtuous', 'Wondrous', 'Adorable', 'Breathtaking', 'Caring', 'Energetic', 'Flourishing', 'Graceful', 'Illuminating', 'Joyous', 'Kinetic', 'Luxurious', 'Opulent', 'Piquant', 'Resplendent', 'Stunning', 'Treasured', 'Verdant', 'Witty', 'Youthful', 'Zestful', 'Affirmative', 'Brilliant', 'Captivating', 'Delicate', 'Exquisite', 'Frisky', 'Gleaming', 'Ineffable', 'Juicy', 'Luminescent', 'Opalescent', 'Pearly', 'Quixotic', 'Ravishing', 'Sprightly', 'Tantalizing', 'Unblemished', 'Voluptuous', 'Winsome', 'Yummy', 'Zealous', 'Ample', 'Blissful', 'Charismatic', 'Divine', 'Ethereal', 'Fragrant', 'Grandiose', 'Heavenly', 'Incomparable', 'Jubilant', 'Kindest', 'Luminous', 'Mellifluous', 'Noble', 'Ornate'];
        $funnyAnimals = ['Honeybee', 'Cuddly Koala', 'Playful Penguin', 'Wise Owl', 'Cheeky Chipmunk', 'Loyal Labrador', 'Graceful Gazelle', 'Charming Cheetah', 'Majestic Mustang', 'Sassy Seahorse', 'Energetic Elephant', 'Spirited Sparrow', 'Dandy Dalmatian', 'Joyful Jaguar', 'Gleeful Gecko', 'Amiable Armadillo', 'Spirited Squirrel', 'Delightful Dolphin', 'Zesty Zebra', 'Grinning Gorilla', 'Jovial Jellyfish', 'Lucky Ladybug', 'Whimsical Wallaby', 'Dynamic Dragonfly', 'Affectionate Alpaca', 'Radiant Rabbit', 'Vivacious Vulture', 'Jubilant Jackrabbit', 'Bouncing Bunny', 'Breezy Butterfly', 'Dynamic Dingo', 'Merry Meerkat', 'Gleaming Goldfish', 'Chirpy Chickadee', 'Friendly Firefly', 'Snuggly Snail', 'Bubbly Bumblebee', 'Witty Woodpecker', 'Marvelous Magpie', 'Enchanting Eel', 'Jolly Jellybean', 'Snazzy Snail', 'Tender Tortoise', 'Merry Mongoose', 'Playful Platypus', 'Wondrous Warthog', 'Gleaming Gopher', 'Dazzling Dachshund', 'Happy Hummingbird', 'Radiant Raccoon', 'Gleeful Gecko', 'Purrfect Panther', 'Dynamic Deer', 'Gleaming Guppy', 'Jubilant Jaguar', 'Chameleon', 'Merry Manatee', 'Gracious Goldfish', 'Playful Pufferfish', 'Jubilant Jellybean', 'Joyful Jackal', 'Gleaming Gibbon', 'Wholesome Wombat', 'Jubilant Jaybird', 'Breezy Bullfinch', 'Delightful Dugong', 'Pleasant Poodle', 'Sincere Starling', 'Gleeful Gnu', 'Delightful Duckling', 'Sprightly Sparrow', 'Jubilant Jerboa', 'Bouncing Butterfly', 'Grinning Gecko', 'Whimsical Wren', 'Marvelous Macaw', 'Jubilant Javelina', 'Jolly Jellybean', 'Gleaming Goldfinch', 'Friendly Fox', 'Wholesome Weasel', 'Merry Martin', 'Radiant Robin', 'Dynamic Dragon', 'Spirited Starfish', 'Daring Dingo', 'Gleaming Gnat', 'Spirited Snail', 'Zesty Zonkey', 'Breezy Bluebird', 'Lovely Lemur', 'Jubilant Jerboa', 'Gleeful Gecko', 'Glorious Grizzly', 'Radiant Rhino', 'Dynamic Dolphin', 'Jolly Jaguar', 'Marvelous Mule', 'Joyful Jellyfish', 'Gleaming Gander'];
        try {
            return $adjectives[random_int(0, count($adjectives) - 1)] . ' ' . $funnyAnimals[random_int(0, count($funnyAnimals) - 1)];
        } catch (Exception $e) {
            return $adjectives[array_rand($adjectives)] . ' ' . $funnyAnimals[array_rand($funnyAnimals)];
        }
    }
}