import React from 'react';
import { Note } from '../../../../api';
import './NotesList.scss';

interface Props {
  notes: Note[];
  activeNoteID?: string;
  onNoteClicked: (id: string) => void;
}

export default class NotesList extends React.Component<Props> {
  render():JSX.Element {
    const notes = this.props.notes.sort((a,b) => {
      return b.created.getTime() - a.created.getTime();
    })

    return(
      <div className='NotesList'>
        { notes.map(this.renderNote.bind(this)) }
        { this.renderNewNote() }
      </div>
    )
  }

  renderNewNote(): JSX.Element {
    const className='row new' + (this.props.activeNoteID === 'new' ? ' active' : '');

    return(
      <div className={className} onClick={() => this.props.onNoteClicked('new')} >
        <p className='title'>New note</p>
      </div>
    );
  }

  renderNote(note: Note): JSX.Element {
    const options = {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: 'numeric', minute: 'numeric',
      hour12: false
    }

    const formattedDate = new Intl.DateTimeFormat('default', options).format(note.created)

    const className = 'row' + (note.id === this.props.activeNoteID ? ' active' : '')

    return(
      <div key={note.id} className={className} onClick={() => this.props.onNoteClicked(note.id)}>
        <p className='title'>{note.title}</p>
        <p className='date'>{formattedDate}</p>
      </div>
    )
  }
}