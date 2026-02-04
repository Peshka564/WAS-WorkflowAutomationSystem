import { type SvgIconComponent } from '@mui/icons-material';

interface UiConfig {
  icon?: SvgIconComponent;
  color: string;
  textColor: string;
}

export const serviceToUi: Record<string, UiConfig> = {
  drive: {
    color: '#71ffc2',
    textColor: '#404040',
  },
  gmail: {
    color: '#e7e7e7',
    textColor: '#404040',
  },
  github: {
    color: '#171717',
    textColor: '#eeeeee',
  },
};
