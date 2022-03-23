import Scoreboard from './scoreboard.js'
import Version from './version.js'
import Gamelist from './gamelist.js'

Vue.createApp({
  components: {
    Scoreboard,
    Version,
    Gamelist
  },
}).mount('#main')