// Wait for the DOM to be fully loaded before attaching event listeners
document.addEventListener('DOMContentLoaded', () => {
    // Get references to the UI elements
    const playSoloBtn = document.getElementById('play-solo');
    const levelSelection = document.getElementById('level-selection');

    playSoloBtn.addEventListener('click', () => {
        levelSelection.classList.toggle('visible');
    });
});
