import { Card, CardContent, Chip, Typography } from '@mui/material';
import type { WorkflowNodeDisplaySelector } from '../types/workflow';
import { serviceToUi } from '../utils/service-to-ui';

export function NodeBase({ data }: { data: WorkflowNodeDisplaySelector }) {
  console.log(data);
  const mainColor = serviceToUi[data.serviceName].color;
  const secondaryColor = serviceToUi[data.serviceName].textColor;

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
          label={data.serviceName}
          sx={{ backgroundColor: secondaryColor, color: mainColor }}
        />
      </CardContent>
    </Card>
  );
}
