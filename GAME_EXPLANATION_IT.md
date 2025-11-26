# Approfondimento - Logica di Simulazione del Gioco

## 🌊 Panoramica
Questo programma simula un gioco da tavolo di sopravvivenza ambientato nelle profondità oceaniche. I giocatori assumono il ruolo di sommozzatori che esplorano l'abisso, affrontando incontri pericolosi, rischi ambientali ed eventi soprannaturali per trovare tesori.

L'obiettivo è raccogliere il maggior numero di **Tesori** e **tornare alla barca** con successo prima di finire l'Ossigeno (O2) o soccombere al Panico.

## ⚙️ Meccaniche Principali

### 1. Ossigeno (O2) come Vita e Risorsa
- **HP**: Il tuo mazzo di carte O2 rappresenta la tua vita. Se finisci le carte da pescare, soffochi.
- **Controllo del Respiro**: All'inizio di ogni turno, *devi* scartare una carta O2 per respirare. Questa carta genera anche un evento/sfida casuale che devi risolvere.
- **Sfidare la Sorte**: Puoi spendere carte O2 per "Esplorare" di più, ma questo ti avvicina alla morte.

### 2. Sistema di Panico 😱
Lo stress è letale. Man mano che fallisci le sfide, il tuo livello di Panico aumenta.
- **Livello 0 (Calmo)**: Tira un **d8** (dado a 8 facce).
- **Livello 1 (Stressato)**: Tira un **d6**.
- **Livello 2 (Nel Panico)**: Tira un **d4**.
- **Livello 3 (Attacco di Panico!)**: Perdi l'intero inventario e 5 carte O2 immediatamente. Il Panico si resetta a 0.

### 3. Riserve Abilità
I giocatori hanno 4 statistiche per superare le sfide:
- ⚔️ **Encounter (Incontri)**: Combattere creature.
- 🌊 **Environment (Ambiente)**: Navigare correnti e pericoli.
- 🔧 **Technical (Tecnica)**: Hackerare o riparare equipaggiamento.
- 🔮 **Supernatural (Soprannaturale)**: Resistere a maledizioni e attacchi mentali.

### 4. Ordine di Turno
- **Round 1**: Determinato da un Tiro di Iniziativa (d8).
- **Round Successivi**: Determinato dal **Punteggio**. Il giocatore che ha ottenuto il risultato migliore (punteggio più alto) nel round precedente agisce per primo.

### 5. Risoluzione delle Sfide (Carte)
Quando peschi una carta O2 (Controllo Respiro o Esplorazione), devi risolverla:
1.  **Difficoltà**: La carta ha un valore numerico (es. Encounter-4).
2.  **Sforzo**: Puoi spendere punti dalla Riserva Abilità corrispondente (es. Encounter) per aggiungerli al tuo tiro.
3.  **Tiro del Dado**: Tira il dado del tuo attuale livello di panico (d8, d6 o d4).
    *   💥 **Dadi Esplosivi**: Se tiri il valore massimo, tira di nuovo e somma il risultato! (Critico)
4.  **Risultato**: `Punti Spesi + Tiro Dado` vs `Difficoltà Carta`.
    *   **Successo (>=)**: Superi la sfida e ottieni **Punteggio** pari al valore della carta.
    *   **Fallimento (<)**: Fallisci e il tuo **Panico aumenta di 1**.

## 🎮 Flusso di Gioco

### Fase 1: Riposo 💤
- Se un giocatore ha **0 Panico**, recupera **1 Punto Abilità** (tipo casuale).
- I giocatori stressati non possono riposare.

### Fase 2: Turni dei Giocatori
Ogni giocatore esegue le seguenti azioni:
1.  **Controllo del Respiro**: Pesca 1 carta O2. Risolvi la sua sfida.
2.  **Azioni** (Fino a 3):
    - **(E)xplore (Esplora)**: Pesca un'altra carta O2, risolvila e cerca di fare punti.
    - **(U)se Item (Usa Oggetto)**: Usa un oggetto dall'inventario (es. Kit Medico, Bombola O2).
    - **(C)alm Down (Calmati)**: Spendi punti abilità per fare una prova di Volontà e ridurre il Panico.
    - **(Q)uit (Lascia)**: Torna alla barca. Sei al sicuro e tieni il tuo tesoro, ma smetti di giocare.
    - **(P)ass (Passa)**: Termina il turno in anticipo.

### Fase 3: Bottino (Draft) 💰
- I giocatori che hanno fatto punti in questo round possono scegliere oggetti dal mercato.
- **Ordine di Scelta**: Chi ha il punteggio più alto sceglie per primo.
- Gli oggetti includono: **Tesori** (Punti Vittoria), **Equipaggiamento** (Bonus statistiche) e **Consumabili**.

## 🏆 Vincere il Gioco
Il gioco termina quando tutti i giocatori sono **Morti** o **Sulla Barca**.
- **Sopravvissuti**: Solo i giocatori sulla barca possono vincere.
- **Condizione di Vittoria**: Chi ha più Tesori vince.
- **Spareggio**: Se i Tesori sono uguali, vince il giocatore che è rimasto in gioco più a lungo (uscito in un round successivo).

## 🛠 Dettagli Tecnici
- **Linguaggio**: Go (Golang).
- **Logging**: Il gioco esporta un log dettagliato turno per turno su `game_logs.xlsx` usando la libreria `excelize`.
- **Simulazione**: Supporta un ciclo `NumberOfGames` per bilanciamento e testing (attualmente impostato a 1).
