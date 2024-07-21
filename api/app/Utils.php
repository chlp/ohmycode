<?php

namespace app;
use DateTime;
use Exception;

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
        return preg_match('/^[0-9a-zA-Z_!?:=+\-,.\s\'АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя]{1,64}$/', $str) === 1;

    }

    static public function randomName(): string
    {
        $adjectives = ['Amiable', 'Blissful', 'Cheerful', 'Delightful', 'Enchanting', 'Friendly', 'Gracious', 'Harmonious', 'Invigorating', 'Jovial', 'Kindhearted', 'Lively', 'Magnificent', 'Nurturing', 'Optimistic', 'Playful', 'Quaint', 'Radiant', 'Serene', 'Tranquil', 'Uplifting', 'Vibrant', 'Wholesome', 'Affectionate', 'Beautiful', 'Charming', 'Dreamy', 'Elegant', 'Festive', 'Gentle', 'Heart warming', 'Inspiring', 'Jubilant', 'Kind', 'Lovely', 'Majestic', 'Noble', 'Outstanding', 'Pleasurable', 'Radiant', 'Splendid', 'Tender', 'Unforgettable', 'Virtuous', 'Wondrous', 'Adorable', 'Breathtaking', 'Caring', 'Energetic', 'Flourishing', 'Graceful', 'Illuminating', 'Joyous', 'Kinetic', 'Luxurious', 'Opulent', 'Piquant', 'Resplendent', 'Stunning', 'Treasured', 'Verdant', 'Witty', 'Youthful', 'Zestful', 'Affirmative', 'Brilliant', 'Captivating', 'Delicate', 'Exquisite', 'Frisky', 'Gleaming', 'Ineffable', 'Juicy', 'Luminescent', 'Opalescent', 'Pearly', 'Quixotic', 'Ravishing', 'Sprightly', 'Tantalizing', 'Unblemished', 'Voluptuous', 'Winsome', 'Yummy', 'Zealous', 'Ample', 'Blissful', 'Charismatic', 'Divine', 'Ethereal', 'Fragrant', 'Grandiose', 'Heavenly', 'Incomparable', 'Jubilant', 'Kindest', 'Luminous', 'Mellifluous', 'Noble', 'Ornate'];
        $animals = ['Honeybee', 'Koala', 'Penguin', 'Owl', 'Chipmunk', 'Labrador', 'Gazelle', 'Cheetah', 'Mustang', 'Seahorse', 'Elephant', 'Sparrow', 'Dalmatian', 'Jaguar', 'Gecko', 'Armadillo', 'Squirrel', 'Dolphin', 'Zebra', 'Gorilla', 'Jellyfish', 'Ladybug', 'Wallaby', 'Dragonfly', 'Alpaca', 'Rabbit', 'Vulture', 'Jackrabbit', 'Bunny', 'Butterfly', 'Dingo', 'Meerkat', 'Goldfish', 'Chickadee', 'Firefly', 'Snail', 'Bumblebee', 'Woodpecker', 'Magpie', 'Eel', 'Jellybean', 'Snail', 'Tortoise', 'Mongoose', 'Platypus', 'Warthog', 'Gopher', 'Dachshund', 'Hummingbird', 'Raccoon', 'Gecko', 'Panther', 'Deer', 'Guppy', 'Jaguar', 'Chameleon', 'Manatee', 'Goldfish', 'Pufferfish', 'Jellybean', 'Jackal', 'Gibbon', 'Wombat', 'Jaybird', 'Bullfinch', 'Dugong', 'Poodle', 'Starling', 'Gnu', 'Duckling', 'Sparrow', 'Jerboa', 'Butterfly', 'Gecko', 'Wren', 'Macaw', 'Javelina', 'Jellybean', 'Goldfinch', 'Fox', 'Weasel', 'Martin', 'Robin', 'Dragon', 'Starfish', 'Dingo', 'Gnat', 'Snail', 'Zonkey', 'Bluebird', 'Lemur', 'Jerboa', 'Gecko', 'Grizzly', 'Rhino', 'Dolphin', 'Jaguar', 'Mule', 'Jellyfish', 'Gander'];
        try {
            return $adjectives[random_int(0, count($adjectives) - 1)] . ' ' . $animals[random_int(0, count($animals) - 1)];
        } catch (Exception $e) {
            return $adjectives[array_rand($adjectives)] . ' ' . $animals[array_rand($animals)];
        }
    }

    static public function timer(): float
    {
        static $start = null;
        if ($start == null) {
            $start = microtime(true);
            return 0;
        }
        return $start;
    }

    static public function log(string $str): void
    {
        $msg = (DateTime::createFromFormat('U.u', microtime(true)))->format('Y-m-d H:i:s.u');
        $timer = self::timer();
        if ($timer !== 0.0) {
            $msg .= ' (' . number_format(microtime(true) - $timer, 3, '.', '') . ')';
        }
        $msg .= ': ';
        $msg .= $str;
        error_log($msg);
    }
}
