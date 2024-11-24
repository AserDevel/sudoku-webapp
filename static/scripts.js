// Variable to keep track of the currently open popup
let currentPopup = null;
let points = 0

// Function to show the number picker
function showNumberPicker(event, row, col) {
    // Prevent form submission or propagation
    event.stopPropagation();

    // Remove the existing popup if one is already open
    if (currentPopup) {
        document.body.removeChild(currentPopup);
        currentPopup = null;
    }

    // Create the popup
    const popup = document.createElement('div');
    popup.classList.add('number-picker');
    popup.style.top = `${event.clientY}px`;
    popup.style.left = `${event.clientX}px`;

    // Add numbers 1-9 to the popup
    for (let i = 1; i <= 9; i++) {
        const btn = document.createElement('button');
        btn.innerText = i;
        btn.onclick = function () {
            // Set the input value when a number is selected
            const cell = document.querySelector(`[name="cell-${row}-${col}"]`);
            if (cell) {
                cell.value = i;
            }
            document.body.removeChild(popup); // Close the popup
            currentPopup = null; 

            sendToServer(row, col, i);
        };
        popup.appendChild(btn);
    }

    // Add the "Erase" button
    const eraseBtn = document.createElement('button');
    eraseBtn.innerText = "X";
    eraseBtn.style.gridColumn = "span 3"; // Make it span across all columns
    eraseBtn.onclick = function () {
        // Clear the input value
        const cell = document.querySelector(`[name="cell-${row}-${col}"]`);
        if (cell) {
            cell.value = "";
        }
        document.body.removeChild(popup); // Close the popup
        currentPopup = null; 

        sendToServer(row, col, 0);
    };
    popup.appendChild(eraseBtn);

    document.body.appendChild(popup);

    // Update the current popup reference
    currentPopup = popup;
}

// Close popups when clicking elsewhere
document.addEventListener('click', function () {
    if (currentPopup) {
        document.body.removeChild(currentPopup);
        currentPopup = null;
    }
});

// Function to send updated cell data to the server
async function sendToServer(row, col, value) {
    // Determine the current difficulty level from the URL path
    const path = window.location.pathname;
    const difficulty = path.split('/')[1]; // Get the first part of the URL path

    // Ensure valid difficulty level is passed
    if (!['easy', 'medium', 'hard'].includes(difficulty)) {
        console.error('Invalid difficulty level');
        return;
    }

    try {
        const response = await fetch('/update-sudoku', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                difficulty: difficulty, // Add the difficulty level to the payload
                row: row,
                col: col,
                value: value, // Send the selected number (or 0 for erase)
            }),
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Server Response:', result);
        } else {
            console.error('Error saving data to the server');
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

// Submits the current sudoku and checks if it's correct
async function submitSudoku(event) {
    event.preventDefault(); // Prevent traditional form submission

    const path = window.location.pathname;
    const difficulty = path.split('/')[1]; // Get the first part of the URL path

    // Ensure valid difficulty level is passed
    if (!['easy', 'medium', 'hard'].includes(difficulty)) {
        console.error('Invalid difficulty level');
        return;
    }

    // Send the JSON with the difficulty, and await response
    try {
        const response = await fetch('/check-sudoku', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ difficulty }), // Send the difficulty only
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();

        // Display the result message
        const messageDiv = document.getElementById('message');
        messageDiv.innerText = result.message;
        messageDiv.style.color = result.correct ? 'green' : 'red';
    } catch (error) {
        console.error('Error:', error);
    }
}

document.getElementById('generate-btn').addEventListener('click', () => {
    const btn = document.getElementById('generate-btn');
    btn.disabled = true; // Disable the button
    btn.innerText = 'Generating...'; // Change the text

    // Reload the page
    const currentPath = window.location.pathname;
    window.location.href = currentPath + "/gen";
});
