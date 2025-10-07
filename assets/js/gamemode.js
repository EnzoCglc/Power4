document.addEventListener('DOMContentLoaded', () => {
    const playSoloBtn = document.getElementById('play-solo');
    const levelSelection = document.getElementById('level-selection');

    playSoloBtn.addEventListener('click', () => {
        levelSelection.classList.toggle('visible');
    });
});
