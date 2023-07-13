import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query';
import { createApp } from 'vue';

import { router } from './router.js';

import App from '@/App.vue';

const app = createApp(App).use(router);

VueQueryPlugin.install(app, {
	queryClient: new QueryClient({
		defaultOptions: {
			queries: {
				refetchOnWindowFocus: false,
				refetchOnMount: false,
				refetchOnReconnect: false,
			},
		},
	}),
});

	app.mount('#app');
