import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { useNavigate, Link } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import {
  Container,
  Grid,
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Box,
  CircularProgress,
  Alert,
  Tooltip,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import EditIcon from '@mui/icons-material/Edit';
import type { Workflow } from '../types/workflow';
import { Pause } from '@mui/icons-material';

const fetchWorkflows = async (): Promise<Workflow[]> => {
  const token = localStorage.getItem('token');
  const response = await axios.get<Workflow[]>(
    'http://localhost:3000/api/workflows',
    {
      headers: { Authorization: `Bearer ${token}` },
    }
  );
  console.log(response.data);
  return response.data;
};

export function Dashboard() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const {
    data: workflows,
    isLoading,
    isError,
    error,
  } = useQuery<Workflow[], Error>({
    queryKey: ['workflows'],
    queryFn: fetchWorkflows,
    refetchOnWindowFocus: true,
  });

  const activateMutation = useMutation({
    mutationFn: async (payload: { id: number; active: boolean }) => {
      const token = localStorage.getItem('token');
      await axios.patch(
        `http://localhost:3000/api/workflows/${payload.id}/activate`,
        {
          active: payload.active,
        },
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workflows'] });
    },
  });

  if (isLoading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="60vh"
      >
        <CircularProgress />
      </Box>
    );
  }

  if (isError) {
    return (
      <Container sx={{ mt: 4 }}>
        <Alert severity="error">
          Error loading workflows:{' '}
          {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
          {(error as any).response?.data?.error || error.message}
        </Alert>
        <Button
          variant="outlined"
          sx={{ mt: 2 }}
          onClick={() => window.location.reload()}
        >
          Retry
        </Button>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={4}
      >
        <Typography variant="h4" component="h1" fontWeight="bold">
          My Workflows
        </Typography>

        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/create-workflow')}
          sx={{ display: { xs: 'none', sm: 'flex' } }}
        >
          Create New
        </Button>
      </Box>

      {workflows?.length === 0 ? (
        <Box
          textAlign="center"
          py={10}
          bgcolor="#f5f5f5"
          borderRadius={2}
          border="1px dashed #ccc"
        >
          <Typography variant="h6" color="textSecondary" gutterBottom>
            No workflows found.
          </Typography>
          <Typography variant="body2" color="textSecondary" mb={3}>
            Get started by creating your first automation!
          </Typography>
          <Button
            variant="contained"
            onClick={() => navigate('/create-workflow')}
          >
            Create Workflow
          </Button>
        </Box>
      ) : (
        <Grid container spacing={3}>
          {(workflows ?? []).map(workflow => (
            <Grid key={workflow.id}>
              <Card
                sx={{
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  transition: 'transform 0.2s',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: 4,
                  },
                }}
              >
                <CardContent sx={{ flexGrow: 1 }}>
                  <Typography gutterBottom variant="h6" component="div" noWrap>
                    {workflow.name}
                  </Typography>

                  <Typography
                    variant="caption"
                    display="block"
                    color="text.secondary"
                    mt={1}
                  >
                    Updated {formatDistanceToNow(new Date(workflow.updated_at))}{' '}
                    ago
                  </Typography>
                  <Typography
                    variant="caption"
                    display="block"
                    color="text.secondary"
                  >
                    Created:{' '}
                    {new Date(workflow.created_at).toLocaleDateString()}
                  </Typography>

                  <Box mt={2} display="flex" alignItems="center" gap={1}>
                    <Box
                      width={10}
                      height={10}
                      borderRadius="50%"
                      bgcolor={
                        workflow.active ? 'success.main' : 'text.disabled'
                      }
                    />
                    <Typography variant="body2" color="text.secondary">
                      {workflow.active ? 'Active' : 'Inactive'}
                    </Typography>
                  </Box>
                </CardContent>

                <CardActions
                  sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}
                >
                  <Button
                    size="small"
                    startIcon={<EditIcon />}
                    component={Link}
                    to={`/workflow/${workflow.id}`}
                  >
                    Edit
                  </Button>

                  <Tooltip title="Trigger Manually">
                    <Button
                      size="small"
                      color="secondary"
                      startIcon={
                        !workflow.active ? <PlayArrowIcon /> : <Pause />
                      }
                      disabled={activateMutation.isPending}
                      onClick={() => {
                        activateMutation.mutate({
                          id: workflow.id,
                          active: !workflow.active,
                        });
                      }}
                    >
                      {workflow.active ? 'Stop' : 'Run'}
                    </Button>
                  </Tooltip>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Container>
  );
}

export default Dashboard;
