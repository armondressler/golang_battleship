export default {
    created() {
        this.fetchData()
    },
    data() {
      return {
        version: null,
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
          let url = `/version`;
          try {
            const version_response = await fetch(url)
            this.error = false;
            if (!version_response.ok) {
              this.error = true;
              this.loading = false;
            }
            this.version = await version_response.json();
            this.loading = false;
          } catch (err) {
            console.log("Failed to fetch version " + err);
            this.loading = false;
            this.error = true;
          }
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
            Failed to fetch version.
        </div>
    </div>
    <div v-else>
        {{ version.version }}
    </div>
    `
}