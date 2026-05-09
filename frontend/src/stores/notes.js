import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '../api/http'

// Notes — записки от руки на доске. Парный класс артефактов рядом с
// reviews/places. См. backend.md §notes и DESIGN-DECISIONS §B2.
export const useNotesStore = defineStore('notes', () => {
  const notes = ref([])
  const myNotes = ref([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      const { data } = await http.get('/notes')
      notes.value = data
    } finally {
      loading.value = false
    }
  }

  async function fetchByAuthor(authorId) {
    const { data } = await http.get('/notes', { params: { author_id: authorId } })
    myNotes.value = data
    return data
  }

  async function create(payload) {
    const { data } = await http.post('/notes', payload)
    notes.value.unshift(data)
    if (data.author_id && myNotes.value.length) myNotes.value.unshift(data)
    return data
  }

  async function update(id, payload) {
    const { data } = await http.put(`/notes/${id}`, payload)
    replaceLocal(data)
    return data
  }

  async function remove(id) {
    await http.delete(`/notes/${id}`)
    notes.value = notes.value.filter(n => n.id !== id)
    myNotes.value = myNotes.value.filter(n => n.id !== id)
  }

  async function strike(id) {
    const { data } = await http.put(`/notes/${id}/strike`)
    replaceLocal(data)
    return data
  }

  function replaceLocal(note) {
    const replaceIn = (arr) => {
      const i = arr.findIndex(n => n.id === note.id)
      if (i !== -1) arr[i] = note
    }
    replaceIn(notes.value)
    replaceIn(myNotes.value)
  }

  return { notes, myNotes, loading, fetchAll, fetchByAuthor, create, update, remove, strike }
})
