import Parser from "rss-parser";
import RSS from "rss";

const parser = new Parser({
    headers: { "User-Agent": "Scaleway Serverless Examples" },
});
const SOURCE_FEED_URL = process.env.SOURCE_FEED_URL || "https://lobste.rs/rss";
const WORTHWILE_TOPICS = process.env.WORTHWILE_TOPICS ||
    "nixos, nix, serverless, terraform";

const feed = new RSS({
    title: `Worthwhile items from ${SOURCE_FEED_URL}`,
    description: `All about ${WORTHWILE_TOPICS}`,
    language: "en",
    ttl: "60",
});

// This is a simple function that given an RSS feed,
// filters out the items that contain any of the topics in the WORTHWILE_TOPICS environment variable.
// It then exposes a RSS feed on its own, with the filtered items.
async function handler(event, _context, _cb) {
    const userAgent = event.headers["User-Agent"];
    const ip = event.headers["X-Forwarded-For"].split(",")[0];
    console.log("Got request from %s with user agent %s", ip, userAgent);

    const topics = WORTHWILE_TOPICS.split(",").map((t) => t.trim());

    const sourceFeed = await parser.parseURL(SOURCE_FEED_URL);
    console.log("Parsing feed: %s", feed.title);

    const worthwhileItems = sourceFeed.items.filter((item) => {
        // Filter out items that don't contain any of the topics
        return topics.some((topic) => item.title.toLowerCase().includes(topic));
    });

    console.log("Found %d worthwhile items", worthwhileItems.length);

    worthwhileItems.forEach((item) => {
        // Shamelessly repost the items
        feed.item({
            title: item.title,
            description: item.content,
            url: item.link,
            date: item.pubDate,
        });
    });

    return {
        statusCode: 200,
        headers: {
            "Content-Type": "application/xml",
        },
        body: feed.xml(),
    };
}

export { handler };

if (process.env.NODE_ENV === "development") {
    import("@scaleway/serverless-functions").then((nodeOffline) => {
        nodeOffline.serveHandler(handler, 8081);
    });
}
