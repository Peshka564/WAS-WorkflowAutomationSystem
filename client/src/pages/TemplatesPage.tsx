import { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Container,
  Paper,
  TextField,
  Typography,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Tooltip,
  Chip,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit'; // Import Edit Icon
import AddIcon from '@mui/icons-material/Add';
import axios from 'axios';

interface Template {
  id: number;
  name: string;
  subject: string;
  body: string;
  email_to: string;
}

const emptyTemplate = { id: 0, name: '', subject: '', body: '', email_to: '' };

export function TemplatesPage() {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [open, setOpen] = useState(false);

  // State now includes ID (0 = New, >0 = Edit)
  const [currentTemplate, setCurrentTemplate] =
    useState<Template>(emptyTemplate);

  useEffect(() => {
    fetchTemplates();
  }, []);

  const fetchTemplates = async () => {
    const token = localStorage.getItem('token');
    try {
      const res = await axios.get('http://localhost:3000/api/templates', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setTemplates(res.data);
    } catch (e) {
      console.error(e);
    }
  };

  // Open modal for Creating (Reset state)
  const handleOpenCreate = () => {
    setCurrentTemplate(emptyTemplate);
    setOpen(true);
  };

  // Open modal for Editing (Load state)
  const handleOpenEdit = (template: Template) => {
    setCurrentTemplate(template);
    setOpen(true);
  };

  // Unified Save Function (Upsert)
  const handleSave = async () => {
    const token = localStorage.getItem('token');
    try {
      await axios.post('http://localhost:3000/api/templates', currentTemplate, {
        headers: { Authorization: `Bearer ${token}` },
      });

      setOpen(false);
      fetchTemplates(); // Refresh list
    } catch (e) {
      console.error('Failed to save template', e);
      alert('Failed to save template');
    }
  };

  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={3}
      >
        <Typography variant="h4" fontWeight="bold">
          Email Templates
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleOpenCreate}
        >
          New Template
        </Button>
      </Box>

      <Paper elevation={2}>
        <List>
          {templates.map(t => (
            <ListItem
              key={t.id}
              divider
              secondaryAction={
                <Box>
                  <Tooltip title="Edit">
                    <IconButton
                      edge="end"
                      onClick={() => handleOpenEdit(t)}
                      sx={{ mr: 1 }}
                    >
                      <EditIcon />
                    </IconButton>
                  </Tooltip>
                </Box>
              }
            >
              <ListItemText
                primary={
                  <Typography variant="subtitle1" fontWeight="medium">
                    {t.name}
                  </Typography>
                }
                secondary={
                  <>
                    <Typography variant="body2" color="text.primary">
                      Subject: {t.subject}
                    </Typography>
                    {t.email_to && (
                      <Chip
                        label={`To: ${t.email_to}`}
                        size="small"
                        variant="outlined"
                        color="primary"
                        sx={{ height: 20, fontSize: '0.7rem' }}
                      />
                    )}
                    <Typography
                      variant="caption"
                      color="textSecondary"
                      noWrap
                      sx={{ maxWidth: 500, display: 'block' }}
                    >
                      {t.body}
                    </Typography>
                  </>
                }
              />
            </ListItem>
          ))}
          {templates.length === 0 && (
            <Box p={4} textAlign="center">
              <Typography color="textSecondary">
                No templates found. Create one to get started!
              </Typography>
            </Box>
          )}
        </List>
      </Paper>

      <Dialog
        open={open}
        onClose={() => setOpen(false)}
        fullWidth
        maxWidth="sm"
      >
        <DialogTitle>
          {currentTemplate.id === 0 ? 'Create New Template' : 'Edit Template'}
        </DialogTitle>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} mt={1}>
            <TextField
              label="Template Name"
              placeholder="e.g. Welcome Email"
              fullWidth
              value={currentTemplate.name}
              onChange={e =>
                setCurrentTemplate({ ...currentTemplate, name: e.target.value })
              }
            />
            <TextField
              label="Subject Line"
              //   placeholder="Hello {{.trigger.name}}!"
              fullWidth
              value={currentTemplate.subject}
              onChange={e =>
                setCurrentTemplate({
                  ...currentTemplate,
                  subject: e.target.value,
                })
              }
              //   helperText="Supports variables like {{.trigger.field}}"
            />

            <TextField
              label="To (Default Recipient)"
              placeholder="e.g. admin@company.com" // or {{.trigger.email}}"
              fullWidth
              value={currentTemplate.email_to}
              onChange={e =>
                setCurrentTemplate({
                  ...currentTemplate,
                  email_to: e.target.value,
                })
              }
              //   helperText="Optional. Can use variables like {{.trigger.email}}"
            />

            <TextField
              label="Email Body"
              multiline
              rows={6}
              fullWidth
              value={currentTemplate.body}
              onChange={e =>
                setCurrentTemplate({ ...currentTemplate, body: e.target.value })
              }
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleSave}>
            {currentTemplate.id === 0 ? 'Create' : 'Update'}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}
