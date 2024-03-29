<script setup lang="ts">
import { ColumnDef, getCoreRowModel, getExpandedRowModel, useVueTable, FlexRender } from '@tanstack/vue-table';
import { type Giveaway } from '@twir/api/messages/giveaways/giveaways';
import { NButton, NSpace, NModal, useThemeVars } from 'naive-ui';
import { ref, h, computed } from 'vue';
import { useI18n } from 'vue-i18n';


import Actions from './actions.vue';

import { useUserAccessFlagChecker } from '@/api/index.js';
import Modal from '@/components/giveaways/modal.vue';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';




const { t } = useI18n();
const themeVars = useThemeVars();

const props = withDefaults(defineProps<{
	giveaways: Giveaway[]
	showHeader?: boolean
	showCreateButton?: boolean
	enableGroups?: boolean,
	showBackground?: boolean
}>(), {
	showHeader: false,
	showCreateButton: false,
	enableGroups: false,
});

const userCanManageGiveaways = useUserAccessFlagChecker('MANAGE_GIVEAWAYS');

const showGiveawayEditModal = ref(false);
const editableGiveaway = ref<Giveaway | null>(null);

function onModalClose() {
	editableGiveaway.value = null;
	showGiveawayEditModal.value = false;
}

const columns: ColumnDef<Giveaway>[] = [
	{
		accessorKey: 'description',
		size: 95,
		header: () => h('div', {}, 'Description'),
		cell: ({ row }) => {
			const chevron = row.getCanExpand();

			return h(
					'div',
					{ class: 'flex gap-2 items-center select-none' },
					[chevron, row.getValue('description') as string],
				);
		},
	},
	{
		accessorKey: 'actions',
		size: 5,
		header: () => h('div', {}, 'Actions'),
		cell: ({ row }) => {
			return h(Actions, {
				row: row.original,
				onEdit: () => {
					editableGiveaway.value = row.original;
					showGiveawayEditModal.value = true;
				},
			});
		},
	},
];

const tableValue = computed(() => props.giveaways);

const table = useVueTable({
	get data() {
		return tableValue.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
	getExpandedRowModel: getExpandedRowModel(),
});

</script>

<template>
	<div>
		<div v-if="showHeader" class="header">
			<div>
				<n-space>
					<n-button
						v-if="showCreateButton"
						secondary
						type="success"
						:disabled="!userCanManageGiveaways"
						@click="() => {
							editableGiveaway = null;
							showGiveawayEditModal = true;
						}"
					>
						{{ t('sharedButtons.create') }}
					</n-button>
				</n-space>
			</div>
		</div>

		<n-modal
			v-model:show="showGiveawayEditModal" :mask-closable="false" :segmented="true" preset="card" :title="editableGiveaway?.description ?? t('giveaways.newGiveawayTitle')" class="modal" :style="{
				width: '1400px',
				height: '90dvh',
			}"
			:on-close="onModalClose"
			content-style="padding: 5px;"
		>
			<modal
				:giveaway="editableGiveaway" @close="() => {
					showGiveawayEditModal = false;
					onModalClose();
				}"
			/>
		</n-modal>

		<div
			class="border rounded-md" :class="{ 'mt-5': showHeader }" :style="{
				backgroundColor: props.showBackground ? themeVars.cardColor : 'inherit',
				color: themeVars.textColor2
			}"
		>
			<Table>
				<TableHeader>
					<TableRow
						v-for="headerGroup in table.getHeaderGroups()"
						:key="headerGroup.id"
						class="border-b"
					>
						<TableHead
							v-for="header in headerGroup.headers"
							:key="header.id"
							:style="{ width: `${header.getSize()}%` }"
						>
							<FlexRender
								v-if="!header.isPlaceholder"
								:render="header.column.columnDef.header"
								:props="header.getContext()"
							/>
						</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<template v-if="table.getRowModel().rows?.length">
						<TableRow
							v-for="row in table.getRowModel().rows" :key="row.id"
							:data-state="row.getIsSelected() ? 'selected' : undefined"
							class="border-b"
							:class="{ 'cursor-pointer': true, }"
						>
							<TableCell
								v-for="cell in row.getVisibleCells()"
								:key="cell.id"
								@click="() => {
									if (row.getCanExpand()) {
										row.getToggleExpandedHandler()()
									}
								}"
							>
								<FlexRender
									:render="cell.column.columnDef.cell"
									:props="cell.getContext()"
								/>
							</TableCell>
						</TableRow>
					</template>
					<template v-else>
						<TableRow>
							<TableCell :colSpan="columns.length" class="h-24 text-center">
								No giveaways
							</TableCell>
						</TableRow>
					</template>
				</TableBody>
			</Table>
		</div>
	</div>
</template>

<style scoped>
.header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	flex-wrap: wrap;
	gap: 8px;
}
</style>
