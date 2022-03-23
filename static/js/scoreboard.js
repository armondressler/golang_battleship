export default {
    props: {
      ranking: Number
    },
    created() {   
        this.fetchData()
    },
    data() {
      return {
        scoreboard: [],
        loading: true,
        error: false,
        loading_template: ``,
        error_template: ``
      }
    },
    methods: {
        async fetchData() {
          this.error = false;
          this.loading = true;
          let url = `/players`;
          if (this.ranking) {
            url = url + "?ranking="+ this.ranking;
          }
          try {
            const scoreboard_response = await fetch(url)
            this.error = false;
            if (!scoreboard_response.ok) {
              this.error = true;
              this.loading = false;
            } else {
              this.scoreboard = await scoreboard_response.json();
            }
          } catch (err) {
            console.log("Failed to fetch scoreboard " + err);
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
          Failed to fetch scoreboard.
        </div>
      </div>

      <table class="table mb-4" v-else>
        <thead>
          <tr>
            <th scope="col">Name</th>
            <th scope="col">Wins</th>
            <th scope="col">Losses</th>
          </tr>
        </thead>
        <div class="alert alert-secondary" role="alert" v-if="Object.keys(scoreboard).length == 0">
          No players found.
        </div>
        <tbody>
          <tr v-for="player in scoreboard">
            <th scope="row">{{ player.name }}</th>
            <td>{{ player.wins }}</td>
            <td>{{ player.losses }}</td>
          </tr>
        </tbody>
      </table>
    `
}