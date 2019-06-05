<template>
  <div id="app">
    <ion-app>
      <ion-split-pane when="sm">
        <ion-menu>
          <ion-content padding>
            <ConfigPanel v-model="config"></ConfigPanel>
          </ion-content>
        </ion-menu>
        <ion-page class="ion-page" main>
          <ion-header mode="md">
            <ion-toolbar mode="md" color="primary">
              <ion-buttons slot="start">
                <ion-menu-toggle>
                  <ion-button>
                    <ion-icon slot="icon-only" name="menu"></ion-icon>
                  </ion-button>
                </ion-menu-toggle>
              </ion-buttons>
              <ion-title>Go slice algorithms code generator</ion-title>
            </ion-toolbar>
          </ion-header>
          <ion-content padding>
            <pre class="line-numbers">
              <code class="language-go" id="goCode">{{ goCode() }}</code>
            </pre>
          </ion-content>
        </ion-page>
      </ion-split-pane>
    </ion-app>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from 'vue-property-decorator';
import { generateSourceCode, GeneratorConfig } from './codegen';

import ConfigPanel from './components/ConfigPanel.vue';

@Component({
  components: {
    ConfigPanel,
  },
})
export default class App extends Vue {
  private config: GeneratorConfig = {
    packageName: 'main',
    funcPrefix: 'MyStruct',
    funcSuffix: '',
    sliceType: '*MyStruct',
    acceptLessThan: true,
    useTimSort: false,
  };

  private goCode(): string {
    return generateSourceCode(this.config);
  }

  private updated() {
    const elem = document.getElementById('goCode') as HTMLElement;
    elem.innerHTML = this.escapeHtml(this.goCode());
    const win = window as any;
    win.Prism.highlightAll();
  }

  private escapeHtml(content: string) {
    return content
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }
}
</script>

<style lang="scss">
body {
  height: 100vh;
}

#app {
  height: 100vh;
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
