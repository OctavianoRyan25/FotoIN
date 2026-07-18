import api from '@/lib/api';

export interface Event {
  id: string;
  photographer_id: string;
  title: string;
  slug: string;
  event_date?: string;
  location?: string;
  description?: string;
  cover_photo_url?: string;
  is_active: boolean;
  expires_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEventInput {
  title: string;
  location?: string;
  description?: string;
  event_date?: string;
}

export interface UpdateEventInput {
  title?: string;
  location?: string;
  description?: string;
  event_date?: string;
  is_active?: boolean;
}

export const eventApi = {
  list: async (page = 1, limit = 20) => {
    const response = await api.get(`/events`, { params: { page, limit } });
    return response.data;
  },

  create: async (data: CreateEventInput) => {
    const response = await api.post('/events', data);
    return response.data;
  },

  getById: async (id: string) => {
    const response = await api.get(`/events/${id}`);
    return response.data;
  },

  update: async (id: string, data: UpdateEventInput) => {
    const response = await api.put(`/events/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    const response = await api.delete(`/events/${id}`);
    return response.data;
  },
};
