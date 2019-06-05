import Vue from 'vue';
import App from './App.vue';
import './registerServiceWorker';

Vue.config.productionTip = false;
Vue.config.ignoredElements = [
  'ion-app',
  'ion-page',
  'ion-header',
  'ion-toolbar',
  'ion-title',
  'ion-content',
  'ion-button',
  'ion-icon',
  'ion-menu-toggle',
  'ion-item',
  'ion-label',
  'ion-checkbox',
  'ion-toggle',
  'ion-item-divider',
  'ion-item-group',
  'ion-buttons',
  'ion-select',
  'ion-select-option',
  'ion-input',
  'ion-list-header',
  'ion-list',
  'ion-split-pane',
  'ion-menu',
];

new Vue({
  render: (h) => h(App),
}).$mount('#app');
