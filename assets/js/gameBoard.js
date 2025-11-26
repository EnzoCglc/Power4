// Game state variables initialized from server
let currentTurn = window.gameState.currentTurn;  // Current player's turn (1 or 2)
let player1 = window.gameState.player1;          // Player 1's username
let gameMode = window.gameState.gameMode || "duo"; // Game mode: "duo" or "bot"
let botLevel = window.gameState.botLevel || 0;   // Bot difficulty level (1-5)
let Finish = false;                              // Whether the game has ended

// Updates the turn indicator display to show whose turn it is
function updatePlayerIndicator(turn) {
    const turnLayer1 = document.querySelector('.turn-hud-layer1');
    const turnLayer2 = document.querySelector('.turn-hud-layer2');

    if (gameMode === "duo") {
        if (turn === 1) {
            turnLayer1.style.display = 'flex';
            turnLayer2.style.display = 'none';
        } else {
            turnLayer1.style.display = 'none';
            turnLayer2.style.display = 'flex';
        }
    }
}

// Initialize the turn indicator display on page load
updatePlayerIndicator(currentTurn);
document.querySelectorAll('.colonne').forEach(col => {
    col.addEventListener('mouseenter', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            if (!cell.classList.contains('black') && !cell.classList.contains('orange')) {
                // Count empty cells to find where piece would land
                let nb_cell = col.querySelectorAll('.cellule').length - col.querySelectorAll('.black, .orange').length;
                let nb_cell_to_highlight = nb_cell - 1;
                let cell_to_hightlight = col.querySelectorAll('.cellule')[nb_cell_to_highlight];

                // Only show hover effect if game is still ongoing
                if (Finish === true) {
                    cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
                } else {
                    // Add hover effect in the current player's color
                    if (currentTurn === 1) {
                        cell_to_hightlight.classList.add('hover-orange', 'hoverNextTurn');
                    } else {
                        cell_to_hightlight.classList.add('hover-black', 'hoverNextTurn');
                    }
                }
            }
        });
    });
    col.addEventListener('mouseleave', () => {
        col.querySelectorAll('.cellule').forEach(cell => {
            cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
        })
    });
});

// Submits a move to the server and processes the response
function playColumn(colIndex) {
    const endpoint = gameMode === "bot" ? '/game/bot/play' : '/game';

    fetch(endpoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            col: colIndex
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                updateGrid(data.data.game);
            } else {
                console.error('Error:', error);
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

// Animates a piece dropping into a cell
function dropToken(colIndex, rowIndex, player) {
    const Col = document.querySelectorAll('.colonne')[colIndex];
    const Cell = Col.querySelectorAll('.cellule')[rowIndex];

    // Clear any existing state
    Cell.classList.remove('black', 'orange', 'animate-drop');

    // Add color and animation based on player
    if (player === 1) {
        Cell.classList.add('orange', 'animate-drop');
    } else {
        Cell.classList.add('black', 'animate-drop');
    }

    // Clean up animation class when animation completes
    Cell.addEventListener('animationEnd', () => {
        Cell.classList.remove('animate-drop');
    })
}

// Updates the visual game board based on server game state
function updateGrid(game) {
    // Update the visual board to match server state
    game.Columns.forEach((col, colIndex) => {
        const Col = document.querySelectorAll('.colonne')[colIndex];

        col.forEach((cell, rowIndex) => {
            const Cell = Col.querySelectorAll('.cellule')[rowIndex];

            // Only animate pieces that are new (not already on the board)
            if (cell === 1 && !Cell.classList.contains('black')) {
                dropToken(colIndex, rowIndex, 1);
            } else if (cell === 2 && !Cell.classList.contains('orange')) {
                dropToken(colIndex, rowIndex, 2);
            }
        });
    });

    // Update local game state variables
    window.gameState.currentTurn = game.CurrentTurn;
    currentTurn = game.CurrentTurn;
    Finish = game.GameOver;

    // Update the turn indicator to show whose turn it is
    updatePlayerIndicator(currentTurn);

    // Clear any hover effects from the board
    document.querySelectorAll('.cellule').forEach(cell => {
        cell.classList.remove('hoverNextTurn', 'hover-black', 'hover-orange');
    });

    // Handle game-over state
    if (Finish === true) {
        const winMsg = document.getElementById('win-msg');

        // Send game result to server for both duo and bot modes
        // Server will decide whether to update ELO based on ranked status
        const body = {
            winner: game.Winner,
            player1: player1,
            player2: "player2",
            isDraw: game.IsDraw
        };

        fetch('/game/result', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        })
            .then(res => res.json())
            .then(data => console.log("Result saved:", data))
            .catch(err => console.error("Error updating ELO:", err));

        // Display appropriate win message based on game mode and winner
        if (game.Winner === 1) {
            winMsg.textContent = `${player1} Win a game`
        } else if (gameMode === "bot") {
            winMsg.textContent = `Bot Level ${botLevel} Win a game`
        } else {
            winMsg.textContent = "Player 2 Win a game"
        }

        // Show the win banner overlay
        document.querySelector('.win-banner-overlay').style.display = 'flex';
    }
}

// Hide the win banner initially (shown only when game ends)
document.querySelector('.win-banner-overlay').style.display = 'none';

const bg = document.getElementById('bg');
function unlockAudio() {
    if (bg.muted) {
        bg.muted = false;
    }
    bg.play().catch(err => console.warn('Audio playback blocked:', err));
}

window.addEventListener('pointerdown', unlockAudio, { once: true });

// Retry the game by reloading the current page with existing query parameters
function retry() {
    const url = new URL(window.location.href);
    const path = url.pathname;
    const params = url.searchParams.toString();

    const retryUrl = params ? `${path}?${params}` : path;

    window.location.href = retryUrl;
}
