<template>
  <ion-list>
    <ion-item-group>
      <ion-list-header>
        <ion-label>Package/File information</ion-label>
      </ion-list-header>
      <ion-item>
        <ion-label position="stacked">Package Name</ion-label>
        <ion-input v-bind:value="getPackageName" @ionInput="onChangePackageName($event)"></ion-input>
      </ion-item>
      <ion-item>
        <ion-label position="stacked">Function Prefix</ion-label>
        <ion-input v-bind:value="getFunctionPrefix" @ionInput="onChangeFunctionPrefix($event)"></ion-input>
      </ion-item>
      <ion-item>
        <ion-label position="stacked">Function Suffix</ion-label>
        <ion-input v-bind:value="getFunctionSuffix" @ionInput="onChangeFunctionSuffix($event)"></ion-input>
      </ion-item>
      </ion-item-group>

      <ion-item-group>
      <ion-list-header>
        <ion-label>Generator Option</ion-label>
      </ion-list-header>
      <ion-item>
        <ion-label>Use Primitive Type</ion-label>
        <ion-toggle color="primary" slot="end" v-bind:checked="doesUsePrimitiveType" @ionChange="onChangeType($event)"></ion-toggle>
     </ion-item>
      <ion-item>
        <ion-label position="stacked" v-bind:color="doesUsePrimitiveType ? 'medium' : undefined">Slice Types</ion-label>
        <ion-input v-bind:value="getSliceType" v-bind:disabled="doesUsePrimitiveType" @ionChange="onChangeSliceType($event)"></ion-input>
      </ion-item>
      <ion-item>
        <ion-label position="stacked" v-bind:color="!doesUsePrimitiveType ? 'medium' : undefined">Primitive Type</ion-label>
        <ion-select v-bind:value="getPrimitiveSliceType" @ionChange="onChangePrimitiveSliceType($event)" okText="Okay" cancelText="Dismiss" v-bind:disabled="!doesUsePrimitiveType">
          <ion-select-option value="int">int</ion-select-option>
          <ion-select-option value="uint">uint</ion-select-option>
          <ion-select-option value="byte">byte</ion-select-option>
          <ion-select-option value="rune">rune</ion-select-option>
          <ion-select-option value="int8">int8</ion-select-option>
          <ion-select-option value="uint8">uint8</ion-select-option>
          <ion-select-option value="int16">int16</ion-select-option>
          <ion-select-option value="uint16">uint16</ion-select-option>
          <ion-select-option value="int32">int32</ion-select-option>
          <ion-select-option value="uint32">uint32</ion-select-option>
          <ion-select-option value="int64">int64</ion-select-option>
          <ion-select-option value="uint64">uint64</ion-select-option>
          <ion-select-option value="float32">float32</ion-select-option>
          <ion-select-option value="float64">float64</ion-select-option>
        </ion-select>
      </ion-item>
      <ion-item>
        <ion-label>Accept LessThan function</ion-label>
        <ion-toggle v-bind:checked="doesAcceptLessThan" color="primary" slot="end" @ionChange="onChangeAcceptLessThan($event)"></ion-toggle>
      </ion-item>
      <ion-item>
        <ion-label>Use TimSort</ion-label>
        <ion-toggle color="primary" slot="end" v-bind:checked="doesUseTimSort" @ionChange="onChangeUseTimSort($event)"></ion-toggle>
      </ion-item>
    </ion-item-group>
  </ion-list>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from 'vue-property-decorator';
import { GeneratorConfig } from '../codegen';

@Component
export default class ConfigPanel extends Vue {
  @Prop() public value!: GeneratorConfig;

  private pakcageName = this.value.packageName;
  private functionPrefix = this.value.funcPrefix;
  private functionSuffix = this.value.funcSuffix;
  private sliceType = this.value.sliceType;
  private primitiveSliceType = 'int';
  private usePrimitiveType = false;
  private acceptLessThan = this.value.acceptLessThan;
  private useTimSort = this.value.useTimSort;

  get getPackageName(): string {
    return this.pakcageName;
  }

  get getFunctionPrefix(): string {
    return this.functionPrefix;
  }

  get getFunctionSuffix(): string {
    return this.functionSuffix;
  }

  get getSliceType(): string {
    return this.sliceType;
  }

  get doesUsePrimitiveType(): boolean {
    return this.usePrimitiveType;
  }

  get getPrimitiveSliceType(): string {
    return this.primitiveSliceType;
  }

  get doesAcceptLessThan(): boolean {
    return this.acceptLessThan;
  }

  get doesUseTimSort(): boolean {
    return this.useTimSort;
  }

  public onChangePackageName(event: any) {
    this.pakcageName = event.target.value;
    this.emitConfig();
  }

  public onChangeFunctionPrefix(event: any) {
    this.functionPrefix = event.target.value;
    this.emitConfig();
  }

  public onChangeFunctionSuffix(event: any) {
    this.functionSuffix = event.target.value;
    this.emitConfig();
  }

  public onChangeType(event: any) {
    this.usePrimitiveType = event.target.checked;
    this.emitConfig();
  }

  public onChangeSliceType(event: any) {
    this.sliceType = event.target.value;
    this.emitConfig();
  }

  public onChangePrimitiveSliceType(event: any) {
    this.primitiveSliceType = event.target.value;
    this.emitConfig();
  }

  public onChangeAcceptLessThan(event: any) {
    this.acceptLessThan = event.target.checked;
    this.emitConfig();
  }

  public onChangeUseTimSort(event: any) {
    this.useTimSort = event.target.checked;
    this.emitConfig();
  }

  @Emit()
  public input(value: GeneratorConfig) {
  }

  private emitConfig() {
    const sliceType = this.usePrimitiveType ? this.primitiveSliceType : this.sliceType;
    const config = {
      packageName: this.pakcageName,
      funcPrefix: this.functionPrefix,
      funcSuffix: this.functionSuffix,
      sliceType,
      acceptLessThan: this.acceptLessThan,
      useTimSort: this.useTimSort,
    };
    this.input(config);
  }
}
</script>