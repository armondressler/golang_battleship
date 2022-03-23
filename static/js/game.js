export default {
    props: {
      gameProperties: Object
    },
    data() {
      return {
        g: {
          id: "",
          participants: [],
          maxparticipants: 0,
          description: "",
          state: 0,
          creation_date: "",
          board_parameters: {
            size_x: 0,
            size_y: 0,
            max_ships: 0,
          },
        },
        loading: false,
        error: false,
      }
    },
    computed: {
      description() {
        return this.gameProperties.description.length > 0 ? this.gameProperties.description : "New Game";
      },
      creationTimeAgo() {
        const now = new Date();
        const game_creation_date = new Date(this.gameProperties.creation_date);
        if (now - game_creation_date < 60 * 1000) {
          return "a few seconds ago";
        } else if (now - game_creation_date <  60 * 60 * 1000) {
          return Math.floor((now - game_creation_date) / (1000 * 60)) + " minutes ago";
        } else if (now - game_creation_date <  24 * 60 * 60 * 1000) {
          return Math.floor((now - game_creation_date) / (1000 * 60 * 60)) + " hours ago";
        } else {
          return Math.floor((now - game_creation_date) / (1000 * 60 * 60 * 24)) + " days ago";
        }
      }
    },
    methods: {
        async fetchData() {
          this.error = false;
          this.loading = true;
          let url = "/games/" + this.gameProperties.id;
          try {
            const game_response = await fetch(url)
            this.error = false;
            if (!game_response.ok) {
              this.error = true;
              this.loading = false;
            } else {
              this.gameProperties = await game_response.json();
            }
          } catch (err) {
            console.log("Failed to fetch game " + err);
            this.error = true;
          }
          this.loading = false;
        },
    },
    template: `
      <div v-if="loading">
        <strong>
          Loading...
        </strong>
        <div class="spinner-border" aria-hidden="true"></div>
      </div>
      <div v-else-if="error">
        <div class="alert alert-danger" role="alert">
          Failed to fetch game.
        </div>
      </div>

      <div class="card">
        <div class="card-body">
          <h5 class="card-title">{{ description }}</h5>
          <p class="card-text">{{ gameProperties.participants.join(", ") }} ( {{ gameProperties.max_participants }} )</p>
          <p class="card-text">
            <small class="text-muted">{{ creationTimeAgo }}</small>
          </p>
          <a href="#" class="btn btn-primary">Join</a>
        </div>
      </div>
    `
}




