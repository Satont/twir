import { createGlobalState, useLocalStorage } from '@vueuse/core'

import type { LayoutItem } from 'grid-layout-plus'

const version = '10'

export type WidgetItem = LayoutItem & { visible: boolean }

export const useWidgets = createGlobalState(() => {
	return useLocalStorage<WidgetItem[]>(`twirWidgetsPositions-v${version}`, [
		{
			x: 6,
			y: 0,
			w: 4,
			h: 8,
			i: 'chat',
			minW: 3,
			minH: 8,
			visible: true,
		},
		{
			x: 3,
			y: 0,
			w: 3,
			h: 8,
			i: 'events',
			minW: 3,
			minH: 8,
			visible: true,
		},
		{
			x: 0,
			y: 0,
			w: 3,
			h: 8,
			i: 'stream',
			minW: 3,
			minH: 8,
			visible: true,
		},
		{
			x: 0,
			y: 9,
			w: 3,
			h: 8,
			i: 'audit-logs',
			minW: 3,
			minH: 8,
			visible: true,
		},
	])
})
