const gameBoard = document.getElementById("gameBoard");
const message = document.getElementById("message");
const columns = 7;
const rows = 6;

const urlParams = new URLSearchParams(window.location.search);
const player1Name = urlParams.get("player1") || "Jugador 1";
const player2Name = urlParams.get("player2") || "Jugador 2";
let currentPlayer = player1Name;

gameBoard.innerHTML = Array.from({ length: rows }, () => Array(columns).fill(0))
	.map(
		(row, rowIndex) =>
			row
				.map(
					(cell, colIndex) =>
						`<div class="cell" data-row="${rowIndex}" data-col="${colIndex}"></div>`
				)
				.join('')
	)
	.join('');

const cells = Array.from(document.querySelectorAll('.cell'));
cells.forEach((cell) => {
	cell.addEventListener('click', handleClick);
});

message.textContent = `${currentPlayer}, es tu turno`;

async function handleClick(event) {
  const col = parseInt(event.target.dataset.col);

  // Hacer una llamada API al backend para realizar un movimiento
  const response = await fetch("/api/make-move?col=" + col, {
    method: "POST",
  });

  const result = await response.json();

  if (result.error) {
    console.error(result.error);
    return;
  }

  const { row, status, player } = result;

  if (row !== null) {
    await animateFallingPiece(row, col, player);

    if (status === "win") {
      message.textContent = `¡${currentPlayer} gana!`;
      cells.forEach((cell) => cell.removeEventListener("click", handleClick));
    } else {
      currentPlayer = currentPlayer === player1Name ? player2Name : player1Name;
      message.textContent = `${currentPlayer}, es tu turno`;
    }
  }
}

async function animateFallingPiece(row, col, player) {
  const cell = document.querySelector(`[data-row="${row}"][data-col="${col}"]`);
  cell.classList.add(player === 1 ? "player1" : "player2");
  cell.classList.add(player === 1 ? "player1-animate" : "player2-animate");
  await new Promise((resolve) => setTimeout(resolve, 100));
  cell.classList.remove(player === 1 ? "player1-animate" : "player2-animate");
}