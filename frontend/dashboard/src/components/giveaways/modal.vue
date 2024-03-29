<script setup lang="ts">


import { type FormInst, type FormItemRule, type FormRules, NForm, NFormItem, NInput, NButton, NSlider, NInputNumber, NModal } from 'naive-ui';
import { computed, onMounted, ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { EditableGiveaway } from './types';


import { useGiveawaysManager } from '@/api';
import { useChooseGiveawayWinners, useParticipants } from '@/api/giveaways';
import Chat from '@/components/giveaways/chat.vue';
import Users from '@/components/giveaways/users.vue';
import Modal from '@/components/giveaways/winners.vue';

const { t } = useI18n();

const props = defineProps<{
	giveaway: EditableGiveaway | null
}>();

const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableGiveaway>({
	description: '',
	rolesIds: [],
	followersAgeLuck: true,
	followersLuck: 0,
	isFinished: false,
	isRunning: false,
	keyword: '',
	requiredMinFollowTime: 0,
	requiredMinMessages: 0,
	requiredMinWatchTime: 0,
	requiredMinSubscriberTier: 0,
	winnersCount: 1,
	requiredMinSubscriberTime: 0,
	subscribersLuck: 0,
});

onMounted(() => {
	if (!props.giveaway) return;

	formValue.value = toRaw(props.giveaway);
});



const giveawaysManager = useGiveawaysManager();
const giveawaysCreate = giveawaysManager.create;
const giveawaysUpdate = giveawaysManager.update;

async function save() {
	await formRef.value?.validate();

	const rawData = toRaw(formValue.value);

	const data = {
		...rawData,
	};

	if (rawData.id) {
		await giveawaysUpdate.mutateAsync({
			...data,
			id: rawData.id,
		});
	} else {
		await giveawaysCreate.mutateAsync({
			...data,
		});
	}

	emits('close');
}

const rules: FormRules = {
	description: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length > 100) {
				return new Error('Description cannot be longer than 100 characters');
			}

			if (!value || !value.length) {
				return new Error('Description is required');
			}

			return true;
		},
	},
	keyword: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length > 100) {
				return new Error('Keyword cannot be longer than 100 characters');
			}

			if (!value || !value.length) {
				return new Error('Keyword is required');
			}
		},
	},
};


const participants = useParticipants(formValue.value.id ?? '', '');
const participantsCount = computed(() => participants.data.value?.totalCount ?? 0);
const isAbleToRoll = computed(
	() => !formValue.value.isFinished && formValue.value.id,
);

const chooseWinners = useChooseGiveawayWinners(formValue.value.id ?? '');

async function onChangeRunningClick() {
	await formRef.value?.validate();
	formValue.value.isRunning = !formValue.value.isRunning;

	const rawData = toRaw(formValue.value);

	await giveawaysUpdate.mutateAsync({
		...rawData,
		id: rawData.id ?? '',
	});
}

async function onFinishClick() {
	await formRef.value?.validate();
	formValue.value.isFinished = true;

	const rawData = toRaw(formValue.value);

	await giveawaysUpdate.mutateAsync({
		...rawData,
		id: rawData.id ?? '',
	});
}

async function onRollClick() {
	if (!isAbleToRoll.value) return;

	await chooseWinners.mutateAsync({
		giveawayId: formValue.value.id ?? '',
	});

	showWinnersModal.value = true;
}

const showWinnersModal = ref(false);

</script>

<template>
	<n-modal
		v-model:show="showWinnersModal" :mask-closable="false" :segmented="true" preset="card" title="Winners" class="modal" :style="{
			width: '500px',
			height: '400px',
		}"
	>
		<Modal :giveaway="giveaway" />
	</n-modal>
	<div class="main-container">
		<div class="flex-container">
			<Users :giveaway="giveaway" @open-winners="showWinnersModal = true" />
		</div>
		<div class="flex-container" style="overflow-y: scroll;">
			<div style="display: flex; flex-direction: column; gap: 12px">
				<n-form ref="formRef" :model="formValue" :rules="rules" class="flex flex-col h-[95%] flex-grow">
					<n-form-item :label="t('giveaways.modal.description')" path="description" show-require-mark style="width: 100%">
						<n-input v-model:value="formValue.description" :maxlength="120" />
					</n-form-item>
					<!-- <n-form-item :label="t('giveaways.modal.permissions')" path="rolesIds">
						<div style="display: flex; flex-direction: column; gap: 5px">
							<n-button-group v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)" :key="index">
								<n-button
									v-for="option of group" :key="option.value" :type="formValue.rolesIds.includes(option.value) ? 'success' : 'default'"
									secondary
									@click="() => {
										if (formValue.rolesIds.includes(option.value)) {
											formValue.rolesIds = formValue.rolesIds.filter(r => r !== option.value)
										} else {
											formValue.rolesIds.push(option.value)
										}
									}"
								>
									<template #icon>
										<IconSquareCheck v-if="formValue.rolesIds.includes(option.value)" />
										<IconSquare v-else />
									</template>
									{{ option.label }}
								</n-button>
							</n-button-group>
						</div>
					</n-form-item> -->
					<n-form-item :label="t('giveaways.modal.keyword')" show-require-mark path="keyword">
						<n-input
							v-model:value="formValue.keyword"
							type="text"
							placeholder="Keyword Phrase"
						/>
					</n-form-item>
					<n-form-item :label="t('giveaways.modal.followersLuck')" path="followersLuck">
						<n-slider v-model:value="formValue.followersLuck" :step="1" :max="10" :min="0" />
					</n-form-item>
					<n-form-item :label="t('giveaways.modal.subscribersLuck')" path="subscribersLuck">
						<n-slider v-model:value="formValue.subscribersLuck" :step="1" :max="10" :min="0" />
					</n-form-item>
					<n-form-item :label="t('giveaways.modal.requiredMinWatchTime')" path="requiredMinWatchTime">
						<n-input-number
							v-model:value="formValue.requiredMinWatchTime"
							placeholder="Minimum watch time"
							min="0"
						/>
					</n-form-item>
					<!-- <n-form-item :label="t('giveaways.modal.requiredMinSubscriberTime')" path="requiredMinSubscriberTime">
						<n-input-number
							v-model:value="formValue.requiredMinSubscriberTime"
							placeholder="Required min subscribe time"
							min="0"
						/>
					</n-form-item> -->
					<n-form-item :label="t('giveaways.modal.requiredMinMessages')" path="requiredMinMessages">
						<n-input-number
							v-model:value="formValue.requiredMinMessages"
							placeholder="Minimum messages"
							min="0"
						/>
					</n-form-item>
					<n-form-item :label="t('giveaways.modal.winnersCount')" path="winnersCount">
						<n-input-number
							v-model:value="formValue.winnersCount"
							placeholder="Winners count"
							min="1"
						/>
					</n-form-item>
				</n-form>
				<n-button style="margin-top: 8px" secondary type="success" block :disabled="formValue.isFinished" @click="save">
					{{ t('sharedButtons.save') }}
				</n-button>
				<div>
					<n-button
						type="primary"
						style="width: 100%; margin-bottom: 5px;"
						:disabled="formValue.isFinished"
						@click="onChangeRunningClick"
					>
						{{
							props.giveaway?.isRunning
								? t("giveaways.modal.running")
								: t("giveaways.modal.stopped")
						}}
					</n-button>
					<n-button
						type="primary"
						style="width: 100%; margin-bottom: 5px;"
						:disabled="formValue.isFinished"
						@click="onFinishClick"
					>
						{{ t("giveaways.modal.finish") }}
					</n-button>
					<n-button
						type="primary"
						style="width: 100%; margin-bottom: 5px;"
						:disabled="!isAbleToRoll"
						@click="onRollClick"
					>
						{{
							isAbleToRoll
								? t("giveaways.modal.roll")
								: t("giveaways.modal.cantRoll")
						}}
					</n-button>
				</div>
			</div>
		</div>
		<div class="flex-container">
			<Chat />
		</div>
	</div>
</template>

<style scoped>
.flex-container {
	flex: 1;
	width: 100%;
	height: 75dvh;
	margin: 10px;
}

.main-container {
	display: flex;
	align-items: center;
	justify-content: center;
	max-width: 100%;
	margin: 0 auto;
}

@media screen and (max-width: 768px) {
	.main-container {
		flex-direction: column;
		margin: 8px 0;
	}
}

</style>
