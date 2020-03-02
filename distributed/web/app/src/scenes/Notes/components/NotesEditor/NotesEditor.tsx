import React from 'react';
import { connect } from 'react-redux';
import Call, { Note } from '../../../../api';
import { updateNote, deleteNote } from '../../../../store/notes';
import bin from '../../../../assets/images/bin.png';
import './NotesEditor.scss';

interface Props {
  history: any;

  note: Note;
  noteID: string;
  autoFocus?: boolean;
  updateNote: (note: Note) => void;
  deleteNote: (note: Note) => void;
}

interface State {
  typingTimer?: NodeJS.Timeout;
}

class NotesEditor extends React.Component<Props, State> {
  readonly state: State = {};

  componentDidUpdate(prevProps: Props) {
    if(!prevProps.note || this.props.noteID === 'new') return;
    if(this.props.note === prevProps.note) return;

    // If there was already a save scheduled, cancel it
    if (this.state.typingTimer) clearTimeout(this.state.typingTimer);

    // Schedule a save for 500ms, enough time for a user to continue
    // typing, extending by another 500ms.
    this.setState({
      typingTimer: setTimeout(this.saveChanges.bind(this), 500),
    });
  }

  saveChanges() {
    console.log("Saving changes to note #", this.props.note.id);
    
    const note = {
      id: this.props.note.id,
      title: this.props.note.title,
      text: this.props.note.text,
    }

    Call('updateNote', { note }).catch(console.warn)
  }

  updateValue(key:string, value: string) {
    let note = { ...this.props.note, [key]: value }

    // If its a new note, create via API then write to redux
    if(this.props.noteID === 'new') {
      Call('createNote', { note: { title: note.title, text: note.text } })
        .then(res => {
          const note = new Note(res.data.note);
          this.props.updateNote(note);
          this.props.history.push('/distributed/notes/' + note.id + '/autoFocus');
        })
        .catch(console.warn)
    } else {
      this.props.updateNote(note);
    }
  }

  onTitleChanged(e: any) {
    this.updateValue('title', e.target.value);
  }

  onTextChanged(e: any) {
    this.updateValue('text', e.target.value);
  }

  onDeleteClicked() {
    // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to delete this note?")) return;

    Call('deleteNote', { note: { id: this.props.noteID } }).catch(console.warn);
    this.props.deleteNote(this.props.note);
    this.props.history.push('/distributed/notes');
  }

  render(): JSX.Element {
    if(!this.props.note) return <form className='NotesEditor' />;

    const { title, text, id } = this.props.note;

    return(
      <form className='NotesEditor'>
        <div className='upper'>
          <input
            type='text'
            value={title}
            autoFocus={this.props.autoFocus || id === 'new'}
            placeholder={id === 'new' ? 'Create a new note' : 'Note title'}
            onChange={this.onTitleChanged.bind(this)} />

          { this.props.noteID === 'new' ? null : 
            <img
              src={bin}
              alt='Delete Note'
              onClick={this.onDeleteClicked.bind(this)} /> }
        </div>

        <textarea
          value={text}
          placeholder='Note text'
          onChange={this.onTextChanged.bind(this)} />
      </form>
    );
  }
}

function mapDispatchToProps(dispatch: Function): any {
  return {
    updateNote: (note: Note) => dispatch(updateNote(note)),
    deleteNote: (note: Note) => dispatch(deleteNote(note)),
  };
}

function mapStateToProps(state: any, ownProps: Props): any {
  return {
    note: state.notes.notes[ownProps.noteID] || new Note({ id: 'new' }),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(NotesEditor);