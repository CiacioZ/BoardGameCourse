# Deep Dive - Game Simulation Logic

## 🌊 Overview
This program simulates a survival board game set in the deep ocean. Players take on the role of divers exploring the abyss, facing dangerous encounters, environmental hazards, and supernatural events to find treasure.

The goal is to collect the most **Treasure** and successfully **return to the boat** before running out of Oxygen (O2) or succumbing to Panic.

## ⚙️ Core Mechanics

### 1. Oxygen (O2) as Life & Resource
- **HP**: Your deck of O2 cards represents your life. If you run out of cards to draw, you suffocate.
- **Breath Check**: At the start of every turn, you *must* discard an O2 card to breathe. This card also generates a random event/challenge you must resolve.
- **Pushing Your Luck**: You can spend O2 cards to "Explore" more, but this brings you closer to death.

### 2. Panic System 😱
Stress is a killer. As you fail challenges, your Panic level rises.
- **Level 0 (Calm)**: Roll a **d8** (8-sided die).
- **Level 1 (Stressed)**: Roll a **d6**.
- **Level 2 (Panicked)**: Roll a **d4**.
- **Level 3 (Panic Attack!)**: You lose your entire inventory and 5 O2 cards immediately. Panic resets to 0.

### 3. Ability Pools
Players have 4 stats to overcome challenges:
- ⚔️ **Encounter**: Fighting creatures.
- 🌊 **Environment**: Navigating currents and hazards.
- 🔧 **Technical**: Hacking or fixing equipment.
- 🔮 **Supernatural**: Resisting curses and mental attacks.

### 4. Turn Order
- **Round 1**: Determined by an Initiative Roll (d8).
- **Next Rounds**: Determined by **Score**. The player who performed best (highest score) in the previous round acts first.

### 5. Resolving Challenges (Card Resolution)
When you draw an O2 card (Breath Check or Explore), you must resolve it:
1.  **Check Difficulty**: The card has a numeric value (e.g., Encounter-4).
2.  **Spend Effort**: You can spend points from the matching Ability Pool (e.g., Encounter) to add to your roll.
3.  **Roll Die**: Roll your current panic die (d8, d6, or d4).
    *   💥 **Exploding Dice**: If you roll the maximum value, roll again and add it! (Crit)
4.  **Result**: `Spent Points + Dice Roll` vs `Card Difficulty`.
    *   **Success (>=)**: You overcome the challenge and gain **Score** equal to the card value.
    *   **Failure (<)**: You fail and your **Panic increases by 1**.

## 🎮 Game Loop

### Phase 1: Resting 💤
- If a player has **0 Panic**, they recover **1 Ability Point** (random type).
- Stressed players cannot rest.

### Phase 2: Player Turns
Each player performs the following:
1.  **Breath Check**: Draw 1 O2 card. Resolve its challenge.
2.  **Actions** (Up to 3):
    - **(E)xplore**: Draw another O2 card, resolve it, and try to score points.
    - **(U)se Item**: Use an item from the inventory (e.g., Medkit, Oxygen Tank).
    - **(C)alm Down**: Spend ability points to make a Willpower check and reduce Panic.
    - **(Q)uit**: Return to the boat. You are safe and keep your treasure, but you stop playing.
    - **(P)ass**: End turn early.

### Phase 3: Loot (Draft) 💰
- Players who scored points this round get to pick items from the market.
- **Draft Order**: Highest score picks first.
- Items include: **Treasure** (Victory Points), **Gear** (Stat boosts), and **Consumables**.

## 🏆 Winning the Game
The game ends when all players are either **Dead** or **On the Boat**.
- **Survivors**: Only players on the boat can win.
- **Victory Condition**: Most Treasure wins.
- **Tie-Breaker**: If Treasure is equal, the player who stayed in the game longer (Quit later) wins.

## 🛠 Technical Details
- **Language**: Go (Golang).
- **Logging**: The game exports a detailed turn-by-turn log to `game_logs.xlsx` using the `excelize` library.
- **Simulation**: Supports `NumberOfGames` loop for balancing and testing (currently set to 1).
