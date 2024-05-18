document.addEventListener("DOMContentLoaded", function () {
    // Отправляем пустой POST запрос при загрузке страницы
    fetch('/album-list', { method: 'POST' })
        .then(response => response.json())
        .then(data => {

            // Получаем данные от сервера и рендерим карточки
            renderArtists(data);
            console.log(data)
            // Добавляем обработчик события click для каждой кнопки "More Info"
            let moreInfoButtons = document.querySelectorAll('.more-info-btn');
            moreInfoButtons.forEach(button => {
                button.addEventListener('click', function () {
                    // Отправляем запрос при нажатии на кнопку "More Info"
                    fetchArtistInfo(this);
                });
            });
        })
        .catch(error => console.error('Error fetching artists data:', error));
});

// Функция для рендеринга карточек артистов
function renderArtists(data) {
    // Получаем контейнер для карточек
    let container = document.querySelector('.cards');

    // Проходимся по данным и создаем карточки
    data.forEach(artist => {
        let card = document.createElement('div');
        card.classList.add('col');
        card.innerHTML = `
            <div class="card rounded mx-auto">
                <img src="${artist.image}" class="card-img-top" alt="${artist.name}">
                <div class="card-body text-center">
                    <h5 class="card-title">${artist.name}</h5>
                    <p class="d-none artist-id">${artist.id}</p>
                    <button class="btn btn-info more-info-btn">More Info</button>
                        <div class="spinner-border text-primary d-none" role="status">
                        <span class="sr-only">Loading...</span>
                        </div>
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}

function fetchArtistInfo(button) {
    // Получаем информацию о соответствующем артисте из данных карточки
    let artistId = button.parentElement.querySelector('.artist-id').innerText;
    var requestData = {
        id: artistId
    }
    console.log(requestData)

    button.style.display = 'none';
    let loader = button.parentElement.querySelector('.spinner-border');
    loader.classList.remove('d-none');


    fetch('/artist-info', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            } else {
                button.style.display = 'inline-block';
                let loader = button.parentElement.querySelector('.spinner-border');
                loader.classList.add('d-none');
            }
            return response.json();
        })
        .then(data => {
            // Получаем имя артиста
            let artistName = button.parentElement.querySelector('.card-title').innerText;
            // Находим модальное окно
            let modalDialog = document.querySelector('.modal-dialog');
            // Заполняем модальное окно данными
            console.log(data)
            let artistInfo = data.ArtistInfo;
            let name = artistInfo.name;
            let creationDate = artistInfo.creationDate;
            let firstAlbum = artistInfo.firstAlbum;

            let artistConcertDates = data.ArtistDates
            let artistConcertLocations = data.ArtistLocations
            let artistRelation = data.ArtistDatesLocations

            modalDialog.innerHTML = `
    <div class="modal-content">
        <div class="modal-header">
            <h5 class="modal-title">${name}</h5>
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>
        <div class="modal-body">
    <p><strong>Creation Date:</strong> ${creationDate}</p>
    <p><strong>First Album:</strong> ${firstAlbum}</p>
    <button class="btn btn-primary detail-button" onclick="toggleDisplay('members')">Members</button>         
    <div id="members" class="d-none">
        <ul>
            ${artistInfo.members.map(member => `<li>${member}</li>`).join('')}
        </ul>
    </div>           
    <button class="btn btn-primary detail-button" onclick="toggleDisplay('locations')">Locations</button>
    <div id="locations" class="d-none">
        <ul>
            ${artistConcertLocations.locations.map(location => `<li>${location}</li>`).join('')}
        </ul>
    </div>

    <button class="btn btn-primary detail-button" onclick="toggleDisplay('dates')">Dates</button>
    <div id="dates" class="d-none">
        <ul>
            ${artistConcertDates.dates.map(date => `<li>${date}</li>`).join('')}
        </ul>
    </div>
    <button class="btn btn-primary detail-button" onclick="toggleDisplay('relations')">Relations</button>
    <div id="relations" class="d-none">
        <ul>
        ${Object.keys(artistRelation.datesLocations).map(location => `
        <li>${location}:
            <ul>
                ${artistRelation.datesLocations[location].map(date => `<li>${date}</li>`).join('')}
            </ul>
        </li>`).join('')}
        </ul>
    </div>
</div>

    </div>
`;
            // Отображаем модальное окно
            $('#exampleModal').modal('show');
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
            alert('Server does not response!!!');
        });

}

function toggleDisplay(id) {
    let element = document.getElementById(id);
    if (element.classList.contains('d-none')) {
        element.classList.remove('d-none');
        element.classList.add('d-block');
    } else if (element.classList.contains('d-block')) {
        element.classList.remove('d-block');
        element.classList.add('d-none');
    }
}


