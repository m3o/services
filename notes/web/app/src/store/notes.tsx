import { Note } from '../api';

// Types
const SET_NOTES = 'notes.set_notes';
const UPDATE_NOTE = 'notes.update_note';
const DELETE_NOTE = 'notes.delete_note';

// Interfaces
interface State {
  notes: Record<string, Note[]>
}

interface Action {
  type: string;
  note?: Note;
  notes?: Note[];
}

// Actions
export function setNotes(notes: Note[]): Action {
  return { type: SET_NOTES, notes };
}

export function updateNote(note: Note): Action {
  return { type: UPDATE_NOTE, note };
}

export function deleteNote(note: Note): Action {
  return { type: DELETE_NOTE, note };
}

// Reducer
export default function(state: State = { notes: {} }, action: Action): State {
  switch(action.type) {
    case SET_NOTES: {
      const notes = action.notes!.reduce((map, n) => ({...map, [n.id]: n}), {});
      return { ...state, notes }
    }
    case UPDATE_NOTE: {
      let notes = Object.keys(state.notes).reduce((map, key) => {
        return {...map, [key]: state.notes[key] };
      }, {});

      notes[action.note!.id] = action.note;
      return { ...state, notes };
    }
    case DELETE_NOTE: {
      let notes = Object.keys(state.notes).reduce((map, key) => {
        if(key === action.note!.id) return map;
        return {...map, [key]: state.notes[key] };
      }, {});

      return { ...state, notes }; 
    }
    default: {
      return state;
    }
  }
}