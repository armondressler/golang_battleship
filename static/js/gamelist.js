import Game from './game.js'

export default {
    components: {
        Game
    },
    created() {   
      this.fetchData()
    },
    data() {
      return {
        navItems: [{text: "Open", state: "open"},
                   {text: "Running", state: "running"},
                   {text: "Finished", state: "finished"}],
        activeNavIndex: 0,
        gamelist: new Object,
        loading: false,
        error: false,
      }
    },
    methods: {
        switchNavItemsActiveState(navIndex) {
          for (const [i,navItem] of this.navItems.entries()) {
            if (i != navIndex) {
              navItem.isActive = false;
              continue;
            }
            navItem.isActive = true;
          }
        },
        async fetchData() {
          this.error = false;
          this.loading = true;
          let url = "/games" + "?state=" + this.navItems[this.activeNavIndex].state
          try {
            const gamelist_response = await fetch(url)
            this.error = false;
            if (!gamelist_response.ok) {
              this.error = true;
              this.loading = false;
            } else {
              this.gamelist = await gamelist_response.json();
            }
          } catch (err) {
            console.log("Failed to fetch games " + err);
            this.error = true;
          }
          this.loading = false;
        },
    },
    template: `
    <div class="d-flex flex-row">
    <ul class="nav nav-pills my-3">
      <li class="nav-item" v-for="(navItem, index) in navItems" :key="index">
        <a class="nav-link" @click="this.activeNavIndex = index; fetchData()" v-bind:class="{ active: index == this.activeNavIndex }" href="javascript:void(0);">{{ navItem.text }}</a>
      </li>   
    </ul>
    <button class="btn btn-outline-secondary ms-auto my-3" type="button" @click="fetchData()">
    <svg fill="currentColor" width="32" height="32" viewBox="0 0 16 16">
      <path fill-rule="evenodd" d="M8 3a5 5 0 1 0 4.546 2.914.5.5 0 0 1 .908-.417A6 6 0 1 1 8 2v1z"/>
      <path d="M8 4.466V.534a.25.25 0 0 1 .41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 0 1 8 4.466z"/>
    </svg>
    </button>
    </div>
    <div v-if="loading">
      <strong>
        Loading...
      </strong>
      <div class="spinner-border" aria-hidden="true"></div>
    </div>
    <div v-else-if="error">
        <div class="alert alert-danger" role="alert">
          Failed to fetch games.
        </div>
    </div>
    <div class="alert alert-secondary" role="alert" v-else-if="Object.keys(gamelist).length == 0">
      No games available at the moment.
    </div>
    <div class="card my-2" v-for="g in gamelist">
        <game :gameProperties=g></game>
    </div>
    `
}