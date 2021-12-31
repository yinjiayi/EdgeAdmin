Vue.component("http-request-scripts-config-box", {
	props: ["vRequestScriptsConfig"],
	data: function () {
		let config = this.vRequestScriptsConfig
		if (config == null) {
			config = {}
		}
		return {
			config: config
		}
	},
	methods: {
		changeInitScript: function (scriptConfig) {
			this.config.onInitScript = scriptConfig
			this.$forceUpdate()
		},
		changeRequestScript: function (scriptConfig) {
			this.config.onRequestScript = scriptConfig
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="requestScriptsJSON" :value="JSON.stringify(config)"/>
	<div class="margin"></div>
	<h4>请求初始化</h4>
	<div>
		<script-config-box id="init-script" :v-script-config="config.onInitScript" comment="在接收到客户端请求之后立即调用。预置req、resp变量。" @change="changeInitScript"></script-config-box>
	</div>
	<h4>准备发送请求</h4>
	<div>
		<script-config-box id="request-script" :v-script-config="config.onRequestScript" comment="在准备好转发客户端请求之前调用。预置req、resp变量。" @change="changeRequestScript"></script-config-box>
	</div>
	<div class="margin"></div>
</div>`
})