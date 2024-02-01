import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';

type DeepRequired<T> = Required<{
	[K in keyof T]: T[K] extends Required<T[K]> ? T[K] : DeepRequired<T[K]>;
}>;

export type DudesSettingsWithOptionalId = DeepRequired<Omit<Settings, 'id'> & { id?: string }>;

export const defaultDudesSettings: DudesSettingsWithOptionalId = {
	id: '',
	dudeSettings: {
		color: '#969696',
		maxLifeTime: 1000 * 60 * 30,
		gravity: 400,
		scale: 4,
		soundsEnabled: true,
		soundsVolume: 0.01,
	},
	messageBoxSettings: {
		enabled: true,
		ignoreCommands: true,
		borderRadius: 10,
		boxColor: '#EEEEEE',
		fontFamily: 'roboto',
		fontSize: 20,
		padding: 10,
		showTime: 5 * 1000,
		fill: '#333333',
	},
	nameBoxSettings: {
		fontFamily: 'roboto',
		fontSize: 18,
		fill: ['#FFFFFF'],
		lineJoin: 'round',
		strokeThickness: 4,
		stroke: '#000000',
		fillGradientStops: [0],
		fillGradientType: 0,
		fontStyle: 'normal',
		fontVariant: 'normal',
		fontWeight: 400,
		dropShadow: false,
		dropShadowAlpha: 1,
		dropShadowAngle: 0,
		dropShadowBlur: 1,
		dropShadowDistance: 1,
		dropShadowColor: '#3AC7D9',
	},
};
