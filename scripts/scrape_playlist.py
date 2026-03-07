"""
Scrape songs from Spotify playlists using the public embed endpoint.
No API key needed - uses the public embed data.

Usage:
    python scripts/scrape_playlist.py <playlist_url> [<playlist_url2> ...]

Output:
    data/songs.txt with format: artist|title|album|year|listen_url
"""
import json
import re
import sys
import time
import urllib.request
import urllib.error
from pathlib import Path


def scrape_spotify_playlist(playlist_url: str) -> list[dict]:
    """Scrape track info from a Spotify playlist using the embed page."""
    match = re.search(r"playlist/([a-zA-Z0-9]+)", playlist_url)
    if not match:
        print(f"  Invalid playlist URL: {playlist_url}", file=sys.stderr)
        return []

    playlist_id = match.group(1)
    embed_url = f"https://open.spotify.com/embed/playlist/{playlist_id}"

    req = urllib.request.Request(embed_url, headers={
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    })

    try:
        with urllib.request.urlopen(req, timeout=15) as resp:
            html = resp.read().decode("utf-8")
    except urllib.error.URLError as e:
        print(f"  Error fetching playlist {playlist_id}: {e}", file=sys.stderr)
        return []

    tracks = []

    json_match = re.search(
        r'<script id="__NEXT_DATA__" type="application/json">(.*?)</script>', html
    )
    if not json_match:
        print("  Could not find __NEXT_DATA__ in embed page", file=sys.stderr)
        return []

    try:
        data = json.loads(json_match.group(1))
        props = data.get("props", {}).get("pageProps", {})
        state = props.get("state", {}).get("data", {}).get("entity", {})
        playlist_name = state.get("name", "")
        track_list = state.get("trackList", [])

        for t in track_list:
            uri = t.get("uri", "")
            track_id = uri.split(":")[-1] if ":" in uri else t.get("uid", "")

            # Artist is in the subtitle field (comma + nbsp separated)
            subtitle = t.get("subtitle", "")
            # Clean up non-breaking spaces and split
            artist = subtitle.replace("\u00a0", " ").split(",")[0].strip()
            if not artist:
                artist = "Unknown"

            tracks.append({
                "artist": artist,
                "title": t.get("title", "Unknown"),
                "album": playlist_name,
                "year": "",
                "url": f"https://open.spotify.com/track/{track_id}",
            })
    except (json.JSONDecodeError, KeyError, AttributeError) as e:
        print(f"  Error parsing JSON: {e}", file=sys.stderr)

    return tracks


def main():
    if len(sys.argv) < 2:
        print("Usage: python scrape_playlist.py <spotify_playlist_url> [...]")
        print("\nYou can also manually create data/songs.txt with format:")
        print("artist|title|album|year|listen_url")
        sys.exit(1)

    all_tracks = []
    seen = set()

    for i, url in enumerate(sys.argv[1:]):
        print(f"[{i+1}/{len(sys.argv)-1}] Scraping: {url}", file=sys.stderr)
        tracks = scrape_spotify_playlist(url)
        added = 0
        for t in tracks:
            key = f"{t['artist']}|{t['title']}"
            if key not in seen:
                seen.add(key)
                all_tracks.append(t)
                added += 1
        print(f"  Found {len(tracks)} tracks ({added} new)", file=sys.stderr)
        # Be nice to Spotify
        if i < len(sys.argv) - 2:
            time.sleep(1)

    if not all_tracks:
        print("\nCould not scrape tracks automatically.", file=sys.stderr)
        sys.exit(1)

    out_path = Path(__file__).parent.parent / "data" / "songs.txt"
    out_path.parent.mkdir(exist_ok=True)

    with open(out_path, "w", encoding="utf-8") as f:
        f.write("# format: artist|title|album|year|listen_url\n")
        f.write(f"# Scraped from {len(sys.argv)-1} playlists - {len(all_tracks)} songs\n")
        for t in all_tracks:
            f.write(f"{t['artist']}|{t['title']}|{t['album']}|{t['year']}|{t['url']}\n")

    print(f"\nWrote {len(all_tracks)} songs to {out_path}", file=sys.stderr)


if __name__ == "__main__":
    main()
