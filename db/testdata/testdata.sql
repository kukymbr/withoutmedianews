INSERT INTO statuses
VALUES (1),
       (2),
       (3);

INSERT INTO tags ("tagId", name, "statusId")
VALUES (1, 'Mascots', 1);

INSERT INTO categories ("categoryId", title, "statusId")
VALUES (1, 'Accidents', 1);

INSERT INTO news (title, "shortText", content, author, "categoryId", "tagIds", "publishedAt",  "statusId")
    VALUES (
            'Drunk cat occurred massive traffic jam in the LA',
            'Breaking news from Los Angeles: a stray cat, apparently intoxicated from spilled alcohol, caused a massive traffic jam yesterday at the busy intersection of 5th and Main.',
            'In an unprecedented incident yesterday, a stray tabby cat believed to be intoxicated by spilled alcohol caused a massive traffic jam on downtown Los Angeles streets. Witnesses reported seeing the feline zigzagging across lanes near the intersection of 5th and Main, prompting drivers to slow down and stop altogether. Authorities suspect the cat may have ingested discarded alcohol from nearby trash cans. Animal control was called to safely retrieve the feline, and traffic was gradually restored after the animal was secured. Experts warn that stray animals consuming alcohol can exhibit unpredictable behavior, posing risks to both themselves and motorists.',
            'Bob the Cat',
            1,
            ARRAY [1],
            '2025-09-15 00:00:00 UTC',
            1
        );