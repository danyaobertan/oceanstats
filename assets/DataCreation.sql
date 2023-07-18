------------------------------------------------------------
-- sensor_groups data creation
------------------------------------------------------------

INSERT INTO sensor_groups (name)
VALUES ('alpha'),
       ('beta'),
       ('gamma'),
       ('delta'),
       ('epsilon'),
       ('zeta'),
       ('eta'),
       ('theta'),
       ('iota'),
       ('kappa'),
       ('lambda'),
       ('mu'),
       ('nu'),
       ('xi'),
       ('omicron'),
       ('pi'),
       ('rho'),
       ('sigma'),
       ('tau'),
       ('upsilon'),
       ('phi'),
       ('chi'),
       ('psi'),
       ('omega');

SELECT *
FROM sensor_groups;


------------------------------------------------------------
-- sensors data creation - provided in initialSql.txt texfile in project
------------------------------------------------------------

SELECT *
FROM sensors;


------------------------------------------------------------
-- fish_species data creation
------------------------------------------------------------

INSERT INTO fish_species (name)
VALUES ('Atlantic Bluefin Tuna'),
       ('Atlantic Cod'),
       ('Atlantic Goliath Grouper'),
       ('Atlantic Salmon'),
       ('Atlantic Trumpetfish'),
       ('Atlantic Wolffish'),
       ('Banded Butterflyfish'),
       ('Beluga Sturgeon'),
       ('Blue Marlin'),
       ('Blue Tang'),
       ('Bluebanded Goby'),
       ('Bluehead Wrasse'),
       ('California Grunion'),
       ('Chilean Common Hake'),
       ('Chilean Jack Mackerel'),
       ('Chinook Salmon'),
       ('Clown Triggerfish'),
       ('Coelacanth'),
       ('Common Clownfish'),
       ('Common Dolphinfish'),
       ('Common Fangtooth'),
       ('Deep Sea Anglerfish'),
       ('Flashlight Fish'),
       ('French Angelfish'),
       ('Great Barracuda'),
       ('Green Moray Eel'),
       ('Guineafowl Puffer'),
       ('John Dory'),
       ('Leafy Seadragon'),
       ('Longsnout Seahorse'),
       ('Mexican Lookdown'),
       ('Nassau Grouper'),
       ('Northern Red Snapper'),
       ('Oarfish'),
       ('Ocean Sunfish'),
       ('Orange Roughy'),
       ('Pacific Blackdragon'),
       ('Pacific Halibut'),
       ('Pacific Herring'),
       ('Pacific Sardine'),
       ('Patagonian Toothfish'),
       ('Peruvian Anchoveta'),
       ('Pink Salmon'),
       ('Pygmy Seahorse'),
       ('Queen Angelfish'),
       ('Queen Parrotfish'),
       ('Red Lionfish'),
       ('Sailfish'),
       ('Sarcastic Fringehead'),
       ('Scarlet Frogfish'),
       ('Scorpionfish'),
       ('Skipjack Tuna'),
       ('Slender Snipe Eel'),
       ('Smalltooth Sawfish'),
       ('Sockeye Salmon'),
       ('Spotted Moray'),
       ('Spotted Porcupinefish'),
       ('Spotted Ratfish'),
       ('Stonefish'),
       ('Stoplight Loosejaw'),
       ('Summer Flounder'),
       ('Swordfish'),
       ('Tan Bristlemouth'),
       ('Threespot Damselfish'),
       ('Tropical Two-Wing Flyingfish'),
       ('Wahoo'),
       ('Whiptail Gulper'),
       ('White-Ring Garden Eel'),
       ('Yellowfin Tuna');

SELECT *
FROM fish_species;