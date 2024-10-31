// Fetch and display all movies on page load
window.onload = function () {
    loadMovies();
};


// Fetch the list of movies from the server
function loadMovies() {
    fetch('http://localhost:8090/movies')  // URL of your Go API
        .then(response => response.json())
        .then(movies => {
            const tableBody = document.querySelector('#moviesTable tbody');
            tableBody.innerHTML = ''; // Clear the table
            movies.forEach(movie => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${movie.Title}</td>
                    <td>${movie.Year}</td>
                    <td>${movie.Rating}</td>
                    <td>${movie.Poster}</td>
                `;
                tableBody.appendChild(row);
            });
        })
        .catch(error => {
            console.error('Error fetching movies:', error);
        });
}

// document.getElementById('addMovieForm').addEventListener('submit', function (e) {
//         e.preventDefault();

//         const movieData = {
//             imdb_id: document.getElementById('imdbId').value,
//             title: document.getElementById('title').value,
//             rating: parseFloat(document.getElementById('rating').value),
//             year: parseInt(document.getElementById('year').value),
//         };

//         fetch('http://localhost:8090/add', {
//             method: 'POST',
//             headers: { 'Content-Type': 'application/json' },
//             body: JSON.stringify(movieData)
//         })
//             .then(response => {
//                 if (response.status === 201) {
//                     alert('Movie added successfully!');
//                     document.getElementById('addMovieForm').reset();
//                     loadMovies(); // Reload the movie list
//                 } else {
//                     alert('Failed to add movie');
//                     document.getElementById('addMovieForm').reset();
//                     loadMovies(); // Reload the movie list
//                 }
//             });
//     });

document.getElementById('addMovieForm').addEventListener('submit', async function(event) {
    event.preventDefault(); // Prevent the default form submission

    // Create a JSON object from the form inputs
    const movieData = {
        imdb_id: document.getElementById('imdb_id').value,
        title: document.getElementById('title').value,
        year: parseInt(document.getElementById('rating').value) || null,
        rating: parseFloat(document.getElementById('year').value) || null,
    };

    try {
        const response = await fetch('/movies', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(movieData)
        });

        // Check if the response is ok
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const responseData = await response.json();
        document.getElementById('addMovieForm').innerHTML = JSON.stringify(responseData, null, 2);
        event.target.reset(); // Reset the form
    } catch (error) {
        document.getElementById('addMovieForm').innerHTML = 'Error: Could not connect to server.';
    }
});

// Event listener for delete form submission
document.getElementById("deleteMovieForm").addEventListener("submit", function (e) {
    e.preventDefault(); // Prevent the form from refreshing the page

    const imdbId = document.getElementById("deleteImdbId").value;

    fetch(`http://localhost:8090/delete/${imdbId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (response.ok) {
            document.getElementById("deleteResult").textContent = "Movie deleted successfully!";
            loadMovies(); // Refresh the list of movies
            document.getElementById("deleteMovieForm").reset(); // Reset the form
        } else {
            response.json().then(data => {
                document.getElementById("deleteResult").textContent = data.message || "Failed to delete movie.";
            });
        }
    })
    .catch(error => {
        console.error("Error deleting movie:", error);
        document.getElementById("deleteResult").textContent = "An error occurred while trying to delete the movie.";
    });
});

// // Event listener for details of movie
// document.getElementById("viewMovieForm").addEventListener("submit", function (e) {
//     e.preventDefault();

//     const imdbId = document.getElementById("viewImdbId").value;

//     fetch(`http://localhost:8090/movies/${imdbId}`)
//         .then(response => {
//             if (!response.ok) {
//                 throw new Error("Movie not found");
//             }
//             return response.json();
//         })
//         .then(movie => {
//             const movieDetailsDiv = document.getElementById("movieDetails");
//             movieDetailsDiv.innerHTML = `
//                 <h3>${movie.Title} (${movie.Year})</h3>
//                 <p><strong>IMDb ID:</strong> ${movie.ImdbId}</p>
//                 <p><strong>Rating:</strong> ${movie.Rating}</p>
//                 <p><strong>Poster:</strong> ${movie.Poster ? `<img src="${movie.Poster}" alt="${movie.Title} Poster" width="100">` : "No poster available"}</p>
//             `;
//         })
//         .catch(error => {
//             document.getElementById("movieDetails").innerHTML = `<p style="color:red;">${error.message}</p>`;
//         });
// });

document.getElementById('fetch-details').addEventListener('click', () => {
    const imdbID = document.getElementById('detail-imdb-id').value.trim();
    const responseDiv = document.getElementById('movie-details-response');

    if (!imdbID) {
        responseDiv.innerHTML = `<div class="alert alert-warning" role="alert">Please enter an IMDb ID.</div>`;
        return;
    }

    responseDiv.innerHTML = createSpinner();

    fetch(`/movies/${imdbID}`)
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
    })
// <-- transform to html element code -->

