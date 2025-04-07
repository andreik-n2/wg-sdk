import { configureWunderGraphOperations } from '@virgograph/sdk';

export default configureWunderGraphOperations({
	operations: {
		defaultConfig: {
			authentication: {
				required: false,
			},
		},
		queries: (config) => ({
			...config,
			liveQuery: {
				enable: true,
				pollingIntervalSeconds: 1,
			},
		}),
		mutations: (config) => ({
			...config,
		}),
		subscriptions: (config) => ({
			...config,
		}),
	},
});
