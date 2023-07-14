import { IconTrash, IconPencil } from '@tabler/icons-vue';
import { DataTableColumns, NButton, NPopconfirm, NSpace, NSwitch, NTag, NText } from 'naive-ui';
import { h, VNode } from 'vue';

import { useCommandsManager } from '@/api/index.js';
import type { ListRowData, EditableCommand } from '@/components/commands/types.js';
import { renderIcon } from '@/helpers/index.js';

type Deleter = ReturnType<typeof useCommandsManager>['deleteOne']
type Patcher = NonNullable<ReturnType<typeof useCommandsManager>['patch']>

export const createListColumns = (
	editCommand: (command: EditableCommand) => void,
	deleter: Deleter,
	patcher: Patcher,
): DataTableColumns<ListRowData> => {
	return [
		{
			title: 'Name',
			key: 'name',
			width: 250,
			render(row) {
				return h(
					NTag,
					{
						bordered: false,
						color: { color: row.isGroup ? row.groupColor : 'rgba(112, 192, 232, 0.16)' },
					},
					{ default: () => row.name },
				);
			},
		},
		{
			title: 'Responses',
			key: 'responses',
			render(row) {
				if (row.isGroup) return;
				return h(
					NText,
					{},
					{
						default: () => {
							if (row.module !== 'CUSTOM') return row.description ?? 'No description';
							return row.responses.length
								? h(NSpace, { vertical: true }, row.responses?.map(r => h('span', null, `${r.text}`)))
								: 'Empty responses';
						},
					},
				);
			},
		},
		{
			title: 'Status',
			key: 'enabled',
			width: 100,
			render(row) {
				if (row.isGroup) return;

				return h(
					NSwitch,
					{
						value: row.enabled,
						onUpdateValue: (value: boolean) => {
							row.enabled = value;
							patcher.mutate({ commandId: row.id, enabled: value });
						},
					},
					{ default: () => row.enabled },
				);
			},
		},
		{
			title: 'Actions',
			key: 'actions',
			width: 150,
			render(row) {
				if (row.isGroup) return;
				const editButton = h(
					NButton,
					{
						type: 'primary',
						size: 'small',
						onClick: () => editCommand(row),
						quaternary: true,
					}, {
						icon: renderIcon(IconPencil),
					});

				const deleteButton = h(
					NPopconfirm,
					{
						onPositiveClick: () => deleter.mutate({ commandId: row.id }),
					},
					{
						trigger: () => h(NButton, {
							type: 'error',
							size: 'small',
							quaternary: true,
						}, {
							default: renderIcon(IconTrash),
						}),
						default: () => 'Are you sure you want to delete this command?',
					},
				);

				const buttons: VNode[] = [editButton];

				if (!row.default) {
					buttons.push(deleteButton);
				}

				return h(NSpace, { }, { default: () => buttons });
			},
		},
	];
};
