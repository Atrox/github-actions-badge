package main

const playgroundHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="GitHub Actions Badge provides you with simple badges for your GitHub READMEs to show your GitHub Actions status.">
    <meta name="author" content="Atrox <hello@atrox.dev>">
    <title>GitHub Actions Badge</title>
    <link rel="stylesheet" href="https://unpkg.com/bulma@0.7.5/css/bulma.min.css">
</head>
<body>
<section class="hero is-dark is-small">
    <div class="hero-body">
        <div class="container has-text-centered">
            <h1 class="title">
                GitHub Actions Badge
            </h1>
        </div>
    </div>
</section>

<section class="section">
    <div class="container">
        <div class="columns">
            <div class="column is-half is-offset-one-quarter">
                <div id="app">

					<div class="field-body">
						<div class="field">
							<label class="label">Repository</label>
							<div class="control">
								<input class="input" type="text" placeholder="atrox/sync-dotenv" v-model="repository">
							</div>
						</div>

						<div class="field">
							<label class="label">Branch</label>
							<div class="control">
								<input class="input" type="text" placeholder="master" v-model="ref">
							</div>
						</div>
					</div>

					
                    <div class="field">
                        <label class="label">Style</label>
                        <div class="control">
                            <div class="select">
                                <select v-model="style">
                                    <option>flat</option>
                                    <option>flat-square</option>
                                    <option>plastic</option>
                                    <option>for-the-badge</option>
                                    <option>popout</option>
                                    <option>popout-square</option>
                                    <option>social</option>
                                </select>
                            </div>
                        </div>
                    </div>

                    <div class="card has-text-centered">
                        <div class="card-content">
                            <a v-bind:href="gotoURL" target="_blank">
                                <img alt="Build Status" v-bind:src="badgeURL"/>
                            </a>

                            <label class="label">Markdown</label>
                            <div class="field has-addons">
                                <div class="control is-expanded">
                                    <input class="input" type="text" readonly="readonly" v-bind:value="markdownSource" v-on:focus="$event.target.select()" ref="markdownInput">
                                </div>
                                <div class="control">
                                    <button class="button is-info" v-on:click="copy($refs.markdownInput)">
                                        Copy
                                    </button>
                                </div>
                            </div>

                            <label class="label">HTML</label>
                            <div class="field has-addons">
                                <div class="control is-expanded">
                                    <input class="input" type="text" readonly="readonly" v-bind:value="htmlSource" v-on:focus="$event.target.select()" ref="htmlInput">
                                </div>
                                <div class="control">
                                    <button class="button is-info" v-on:click="copy($refs.htmlInput)">
                                        Copy
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>

<footer class="footer">
    <div class="content has-text-centered">
        <p>
            <strong>GitHub Actions Badge</strong> by <a href="https://atrox.dev">atrox</a><br>
            The source code is available on <a href="https://github.com/atrox/github-actions-badge"><i class="fab fa-github"></i> GitHub</a>
        </p>
    </div>
</footer>

<script src="https://unpkg.com/vue@2.6.10/dist/vue.min.js"></script>
<script>
  var app = new Vue({
    el: '#app',
    data: {
      repository: '',
      style: 'flat',
	  ref: '',
    },
    computed: {
      badgeURL: function () {
        var repo = this.repository || 'atrox/sync-dotenv'
        var style = this.style || ''
        
        var ref = this.ref ? '?ref=' + this.ref : ''
        var url = encodeURIComponent('https://actions-badge.atrox.dev/' + repo + '/badge' + ref)

        return 'https://img.shields.io/endpoint.svg?url=' + url + '&style=' + style
      },
      gotoURL: function () {
        var repo = this.repository || 'atrox/sync-dotenv'
        var ref = this.ref ? '?ref=' + this.ref : ''

        return 'https://actions-badge.atrox.dev/' + repo + '/goto' + ref
      },
      markdownSource: function () {
        return '[![Build Status](' + this.badgeURL + ')](' + this.gotoURL + ')'
      },
      htmlSource: function () {
        return '<a href="' + this.gotoURL + '">' +
          '<img alt="Build Status" src="' + this.badgeURL + '" />' +
          '</a>'
      }
    },
    methods: {
      copy (element) {
        element.select()

        try {
          document.execCommand('copy')
        } catch (err) {
          alert('Oops, unable to copy')
        }
      },
    }
  })
</script>
<script defer src="https://use.fontawesome.com/releases/v5.8.2/js/brands.js" integrity="sha384-GtvEzzhN52RvAD7CnSR7TcPw555abR8NK28tAqa/GgIDk59o0TsaK6FHglLstzCf" crossorigin="anonymous"></script>
<script defer src="https://use.fontawesome.com/releases/v5.8.2/js/fontawesome.js" integrity="sha384-Ia7KZbX22R7DDSbxNmxHqPQ15ceNzg2U4h5A8dy3K47G2fV1k658BTxXjp7rdhXa" crossorigin="anonymous"></script>
</body>
</html>
`
