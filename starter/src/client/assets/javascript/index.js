let store = {
    track_id: undefined,
    track_name: undefined,
    player_id: undefined,
    player_name: undefined,
    race_id: undefined,
}

document.addEventListener("DOMContentLoaded", function () {
    onPageLoad()
    setupClickHandlers()
})

async function onPageLoad() {
    console.log("Getting form info for dropdowns!")
    try {
        getTracks()
          .then(tracks => {
              const html = renderTrackCards(tracks)
              renderAt('#tracks', html)
          })

        getRacers()
          .then((racers) => {
              const html = renderRacerCars(racers)
              renderAt('#racers', html)
          })
    } catch (error) {
        console.log(error)
    }
}

function setupClickHandlers() {
    document.addEventListener('click', function (event) {
        const {target} = event

        // Race track form field
        if (target.matches('.card.track')) {
            handleSelectTrack(target)
            store.track_id = target.id
            store.track_name = target.innerHTML
        }

        // Racer form field
        if (target.matches('.card.racer')) {
            handleSelectRacer(target)
            store.player_id = target.id
            store.player_name = target.innerHTML
        }

        // Submit create race form
        if (target.matches('#submit-create-race')) {
            event.preventDefault()

            // start race
            handleCreateRace()
        }

        // Handle acceleration click
        if (target.matches('#gas-peddle')) {
            handleAccelerate()
        }

        console.log("Store updated :: ", store)
    }, false)
}

async function delay(ms) {
    try {
        return await new Promise(resolve => setTimeout(resolve, ms));
    } catch (error) {
        console.log("an error shouldn't be possible here")
        console.log(error)
    }
}

async function handleCreateRace() {

    const playerId = store.player_id;
    const trackId = store.track_id;
    try {
        const race = await createRace(playerId, trackId);
        store.race_id = race.ID;
        renderAt('#race', renderRaceStartView(store.track_name))
    } catch (error) {
        console.log(error);
    }

    await runCountdown();
    await startRace(store.race_id);
    await runRace(store.race_id);
}
function render() {
    resultsView(race.positions)
}
async function runRace(raceID) {
    return new Promise(resolve => {
        const raceInterval = setInterval(async () => {
            const race = await getRace(raceID);
            console.log('race.positions',race.positions)
            console.log('race.status',race.status)
            console.log('race.status check',(race.status == 'in-progress'))
            console.log('race.status check',(race.status == 'finished'))
            if (race.status == 'in-progress') {
                renderAt('#leaderBoard', raceProgress(race.positions))
            } else if (race.status == 'finished') {
                clearInterval(raceInterval) // to stop the interval from repeating
                renderAt('#race', resultsView(race.positions)) // to render the results view
                resolve(race) // resolve the promise
            }
        }, 500);
    }).catch(error => console.log(error));
}

async function runCountdown() {
    try {
        // wait for the DOM to load
        await delay(1000)
        let timer = 3

        return new Promise(resolve => {
            const countdownInterval = setInterval(() => {
                if (timer !== 0) {
                    document.getElementById('big-numbers').innerHTML = --timer;
                } else {
                    clearInterval(countdownInterval);
                    resolve();
                }
            }, 1000);
        })
    } catch (error) {
        console.log(error);
    }
}

function handleSelectRacer(target) {
    console.log("selected a racer", target.id)

    // remove class selected from all racer options
    const selected = document.querySelector('#racers .selected')
    if (selected) {
        selected.classList.remove('selected')
    }

    // add class selected to current target
    target.classList.add('selected')
}

function handleSelectTrack(target) {
    console.log("selected track", target.id)

    // remove class selected from all track options
    const selected = document.querySelector('#tracks .selected')
    if (selected) {
        selected.classList.remove('selected')
    }

    // add class selected to current target
    target.classList.add('selected')
}

function handleAccelerate() {
    console.log("accelerate button clicked")
    accelerate(store.race_id);
}

function renderRacerCars(racers) {
    console.log(typeof racers)
    console.log(Array.isArray(racers))
    if (!racers.length) {
        return `
			<h4>Loading Racers...</4>
		`
    }

    const results = racers.map(renderRacerCard).join('')

    return `
		<ul id="racers">
			${results}
		</ul>
	`
}

function renderRacerCard(racer) {
    const {id, driver_name, top_speed, acceleration, handling} = racer
    // OPTIONAL: There is more data given about the race cars than we use in the game, if you want to factor in top speed, acceleration,
    // and handling to the various vehicles, it is already provided by the API!
    return `<h4 class="card racer" id="${id}">${driver_name}</h3>`
}

function renderTrackCards(tracks) {
    if (!tracks.length) {
        return `
			<h4>Loading Tracks...</4>
		`
    }

    const results = tracks.map(renderTrackCard).join('')

    return `
		<ul id="tracks">
			${results}
		</ul>
	`
}

function renderTrackCard(track) {
    const {id, name} = track

    return `<h4 id="${id}" class="card track">${name}</h4>`
}

function renderCountdown(count) {
    return `
		<h2>Race Starts In...</h2>
		<p id="big-numbers">${count}</p>
	`
}

function renderRaceStartView(track) {
    return `

<button type="button" onclick="render()">render</button>
		<header>
			<h1>Race: ${track}</h1>
		</header>
		<main id="two-columns">
			<section id="leaderBoard">
				${renderCountdown(3)}
			</section>

			<section id="accelerate">
				<h2>Directions</h2>
				<p>Click the button as fast as you can to make your racer go faster!</p>
				<button id="gas-peddle">Click Me To Win!</button>
			</section>
		</main>
		<footer></footer>
	`
}

function resultsView(positions) {
    positions.sort((a, b) => (a.final_position > b.final_position ? 1 : -1));

    return `
		<header>
			<h1>Race Results</h1>
		</header>
		<main class="result">
			<h3>Race Results</h3>
			<p>The race is done! Here are the final results:</p>
			${raceProgress(positions)}
			<a href="/race">Start a new race</a>
		</main>
	`
}

function raceProgress(positions) {
    let userPlayer = positions.find(e => e.id === parseInt(store.player_id))
    userPlayer.driver_name += " (you)"

    positions = positions.sort((a, b) => (a.segment > b.segment) ? -1 : 1)
    let count = 1

    const results = positions.map(p => {
        return `
			<tr>
				<td>
					<h3>${count++} - ${p.driver_name}</h3>
				</td>
			</tr>
		`
    })

    return `
		<table>
			${results.join('')}
		</table>
	`
}

function renderAt(element, html) {
    const node = document.querySelector(element)

    node.innerHTML = html
}

// ^ Provided code ^ do not remove

const SERVER = 'http://localhost:3001'

function defaultFetchOpts() {
    return {
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': SERVER,
        },
    }
}

async function getTracks() {
    console.log(`calling server :: ${SERVER}/api/tracks`)

    try {
        const data = await fetch(`${SERVER}/api/tracks`)
          .then(response => response.json());
        console.log('tracks', data);
        return data;
    } catch (error) {
        console.log("Problem with getTracks request::", error)
    }
}

async function getRacers() {
    console.log(`calling server :: ${SERVER}/api/cars`)

    try {
        const data = await fetch(`${SERVER}/api/cars`)
          .then(response => response.json());
        console.log('tracks', data);
        return data;
    } catch (error) {
        console.log("Problem with getRacers request::", error)
    }
}

async function createRace(player_id, track_id) {
    player_id = parseInt(player_id)
    track_id = parseInt(track_id)
    const body = {player_id, track_id}

    try {
        const data = await fetch(`${SERVER}/api/races`, {
            method: 'POST',
            ...defaultFetchOpts(),
            dataType: 'jsonp',
            body: JSON.stringify(body)
        }).then(response => response.json())
        console.log('race', data);
        return data;
    } catch (error) {
        console.log("Problem with createRace request::", error)
    }
}

async function getRace(id) {
    console.log(`calling server :: ${SERVER}/api/races/${id}`)

    try {
        const data = await fetch(`${SERVER}/api/races/${id}`)
          .then(response => response.json());
        console.log('tracks', data);
        return data;
    } catch (error) {
        console.log("Problem with getRace request::", error)
    }
}

function startRace(id) {
    console.log(`calling server :: ${SERVER}/api/races/${id}/start`)

    return fetch(`${SERVER}/api/races/${id}/start`, {
        method: 'POST',
        ...defaultFetchOpts(),
    }).catch(error => console.log("Problem with startRace request::", error))
}

async function accelerate(id) {
    console.log(`calling server :: ${SERVER}/api/races/${id}/accelerate`)
    console.log('id', id);

    await fetch(`${SERVER}/api/races/${id}/accelerate`, {
        method: 'POST',
        ...defaultFetchOpts(),
    }).catch(error => console.log("Problem with accelerate request::", error))
}
