#!/usr/bin/env python3

import json
import sqlite3
import sys
import os
from datetime import datetime
import requests
from typing import List, Dict, Any, Optional

ANILIST_API = "https://graphql.anilist.co"

SEASON_QUERY = """
query ($page: Int, $perPage: Int, $season: MediaSeason, $seasonYear: Int) {
  Page(page: $page, perPage: $perPage) {
    media(type: ANIME, season: $season, seasonYear: $seasonYear, sort: POPULARITY_DESC, isAdult: false) {
      id
      title {
        romaji
        english
      }
      description
      coverImage {
        large
      }
      bannerImage
      status
      format
      episodes
      duration
      season
      seasonYear
      genres
      averageScore
      popularity
      tags {
        name
        category
      }
    }
  }
}
"""

def is_hentai(anime: Dict[str, Any]) -> bool:
    return False

def get_seasons_to_fetch() -> List[tuple]:
    current_year = datetime.now().year
    current_month = datetime.now().month
    
    if current_month in [12, 1, 2]:
        current_season = "WINTER"
    elif current_month in [3, 4, 5]:
        current_season = "SPRING"
    elif current_month in [6, 7, 8]:
        current_season = "SUMMER"
    else:
        current_season = "FALL"
    
    seasons = []
    seasons_order = ["WINTER", "SPRING", "SUMMER", "FALL"]
    
    for year in range(current_year + 2, current_year - 3, -1):
        for season in seasons_order:
            seasons.append((season, year))
    
    return seasons

def fetch_anime_data(query: str, variables: Dict[str, Any], max_pages: int = 5) -> List[Dict[str, Any]]:
    all_media = []
    
    for page in range(1, max_pages + 1):
        variables["page"] = page
        variables["perPage"] = 50
        
        try:
            response = requests.post(
                ANILIST_API,
                json={"query": query, "variables": variables},
                headers={"Content-Type": "application/json"}
            )
            response.raise_for_status()
            
            data = response.json()
            if "errors" in data:
                print(f"GraphQL errors on page {page}: {data['errors']}")
                continue
                
            media_list = data["data"]["Page"]["media"]
            if not media_list:
                break
                
            all_media.extend(media_list)
            
            print(f"Fetched {len(media_list)} anime from page {page}")
            
        except requests.exceptions.RequestException as e:
            print(f"Error fetching page {page}: {e}")
            continue
    
    return all_media

def clean_description(description: str) -> str:
    if not description:
        return ""
    
    import re
    clean = re.compile('<.*?>')
    return re.sub(clean, '', description).strip()

def insert_anime_data(db_path: str, anime_list: List[Dict[str, Any]]) -> None:
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    inserted_count = 0
    updated_count = 0
    
    for anime in anime_list:
        try:
            anilist_id = anime["id"]
            title_romaji = anime["title"]["romaji"]
            title_english = anime["title"].get("english")
            description = clean_description(anime.get("description", ""))
            cover_image = anime.get("coverImage", {}).get("large")
            banner_image = anime.get("bannerImage")
            status = anime.get("status")
            format_type = anime.get("format")
            episodes = anime.get("episodes")
            duration = anime.get("duration")
            season = anime.get("season")
            season_year = anime.get("seasonYear")
            genres = ", ".join(anime.get("genres", [])) if anime.get("genres") else None
            score = anime.get("averageScore")
            popularity = anime.get("popularity")
            
            title = title_romaji or title_english or "Unknown Title"
            
            cursor.execute("SELECT id FROM anime WHERE anilist_id = ?", (anilist_id,))
            existing = cursor.fetchone()
            
            if existing:
                cursor.execute("""
                    UPDATE anime SET 
                        title = ?, title_english = ?, title_romaji = ?, description = ?,
                        cover_image = ?, banner_image = ?, status = ?, format = ?,
                        episodes = ?, duration = ?, season = ?, season_year = ?, genres = ?,
                        score = ?, popularity = ?, updated_at = CURRENT_TIMESTAMP
                    WHERE anilist_id = ?
                """, (
                    title, title_english, title_romaji, description,
                    cover_image, banner_image, status, format_type,
                    episodes, duration, season, season_year, genres,
                    score, popularity, anilist_id
                ))
                updated_count += 1
            else:
                cursor.execute("""
                    INSERT INTO anime (
                        anilist_id, title, title_english, title_romaji, description,
                        cover_image, banner_image, status, format, episodes, duration,
                        season, season_year, genres, score, popularity
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    anilist_id, title, title_english, title_romaji, description,
                    cover_image, banner_image, status, format_type, episodes, duration,
                    season, season_year, genres, score, popularity
                ))
                inserted_count += 1
                
        except Exception as e:
            print(f"Error processing anime {anime.get('id', 'unknown')}: {e}")
            continue
    
    conn.commit()
    conn.close()
    
    print(f"Database update complete: {inserted_count} new anime, {updated_count} updated")

def main():
    db_path = os.getenv("DB_PATH", "./anime.db")
    
    print("Starting anime data sync by season...")
    print(f"Database path: {db_path}")
    
    seasons = get_seasons_to_fetch()
    print(f"Fetching {len(seasons)} seasons...")
    
    all_anime = []
    total_fetched = 0
    
    for season, year in seasons:
        print(f"\nFetching {season} {year}...")
        
        variables = {
            "season": season,
            "seasonYear": year
        }
        
        season_anime = fetch_anime_data(SEASON_QUERY, variables, max_pages=3)
        
        if season_anime:
            print(f"Found {len(season_anime)} anime for {season} {year}")
            all_anime.extend(season_anime)
            total_fetched += len(season_anime)
        else:
            print(f"No anime found for {season} {year}")
    
    print(f"\nTotal anime fetched: {total_fetched}")
    
    print("\nStoring data in database...")
    insert_anime_data(db_path, all_anime)
    
    print("Anime data sync completed successfully!")

if __name__ == "__main__":
    main() 