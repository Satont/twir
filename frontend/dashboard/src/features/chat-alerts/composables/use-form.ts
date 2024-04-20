import { useDebounceFn, watchDebounced } from '@vueuse/core';
import { useNotification } from 'naive-ui';
import { defineStore, storeToRefs } from 'pinia';
import type { KeysOfUnion, RequiredDeep, SetNonNullable } from 'type-fest';
import { ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { ChatAlerts, useChatAlertsApi } from '@/api/chat-alerts.js';
import { ChatAlertsSettings } from '@/gql/graphql';

export type FormKey = Exclude<KeysOfUnion<RequiredDeep<SetNonNullable<ChatAlerts>>>, '__typename'>
type Form = Record<FormKey, ChatAlertsSettings>

export const useForm = defineStore('chat-alerts/form', () => {
	const message = useNotification();
	const { t } = useI18n();
	const formRef = ref<HTMLFormElement>();

	const chatAlertsApi = useChatAlertsApi();
	const updateChatAlerts = chatAlertsApi.useMutationUpdateChatAlerts();
	const { chatAlerts: data } = storeToRefs(chatAlertsApi);

	const formValue = ref<Form>({
		chatCleared: {
			enabled: false,
			messages: [],
			cooldown: 2,
		},
		cheers: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		donations: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		firstUserMessage: {
			enabled: false,
			messages: [],
			cooldown: 2,
		},
		followers: {
			enabled: false,
			messages: [],
			cooldown: 3,
		},
		raids: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		redemptions: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		streamOffline: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		streamOnline: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		subscribers: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		ban: {
			enabled: false,
			messages: [],
			cooldown: 2,
			ignoreTimeoutFrom: [],
		},
		unbanRequestCreate: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		unbanRequestResolve: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
	});

	const formInited = ref(false);

	watch(data, (v) => {
		if (!v || formInited.value) return;
		formInited.value = false;

		for (const key of Object.keys(formValue.value) as FormKey[]) {
			if (!v[key]) continue;
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-expect-error
			formValue.value[key] = v[key];
		}

		formInited.value = true;
	}, { immediate: true });

	watchDebounced(formValue, () => {
		if (!formInited.value) {
			return;
		}

		if (!formRef.value) return;
		if (!formRef.value?.reportValidity()) return;

		save();
	}, { deep: true, debounce: 500 });

	async function save() {
		const input = toRaw(formValue.value);
		if (!input) return;

		try {
			await updateChatAlerts.executeMutation({ input });
		} catch (error) {
			message.error({
				title: t('sharedTexts.errorOnSave'),
				duration: 2500,
			});
		}
	}

	const debouncedSave = useDebounceFn(save, 500);

	return { formValue, formInited, save: debouncedSave, formRef };
});
