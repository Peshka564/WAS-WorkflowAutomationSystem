import { Card, CardContent, Chip, Typography } from '@mui/material';
import type { WorkflowNodeDisplaySelector } from '../types/workflow';
import { serviceToUi } from '../utils/service-to-ui';

export function NodeBase({ data }: { data: WorkflowNodeDisplaySelector }) {
  const mainColor = serviceToUi[data.service].color;
  const secondaryColor = serviceToUi[data.service].textColor;

  return (
    <Card
      variant="outlined"
      sx={{
        borderRadius: '20px',
        width: '200px',
        backgroundColor: mainColor,
      }}
    >
      <CardContent
        sx={{
          display: 'flex',
          flexDirection: 'row',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}
      >
        <Typography fontSize={14} color={secondaryColor}>
          {data.taskName}
        </Typography>
        <Chip
          label={data.service}
          sx={{ backgroundColor: secondaryColor, color: mainColor }}
        />
      </CardContent>
    </Card>
  );
}
