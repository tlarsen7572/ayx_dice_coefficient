import AlteryxPythonSDK as Sdk
import xml.etree.ElementTree as Et


class AyxPlugin:
    def __init__(self, n_tool_id: int, alteryx_engine: object, output_anchor_mgr: object):
        # Default properties
        self.n_tool_id: int = n_tool_id
        self.alteryx_engine: Sdk.AlteryxEngine = alteryx_engine
        self.output_anchor_mgr: Sdk.OutputAnchorManager = output_anchor_mgr
        self.label = "Dice Coefficient"

        # Custom properties
        self.Output: Sdk.OutputAnchor = None
        self.OutputField: str = None
        self.Text1: str = None
        self.Text2: str = None

    def pi_init(self, str_xml: str):
        xml_parser = Et.fromstring(str_xml)
        self.OutputField = xml_parser.find("OutputField").text if 'OutputField' in str_xml else ''
        self.Text1 = xml_parser.find("Text1").text if 'Text1' in str_xml else ''
        self.Text2 = xml_parser.find("Text2").text if 'Text2' in str_xml else ''

        # Getting the output anchor from Config.xml by the output connection name
        self.Output = self.output_anchor_mgr.get_output_anchor('Output')

    def pi_add_incoming_connection(self, str_type: str, str_name: str) -> object:
        return IncomingInterface(self)

    def pi_add_outgoing_connection(self, str_name: str) -> bool:
        return True

    def pi_push_all_records(self, n_record_limit: int) -> bool:
        return False

    def pi_close(self, b_has_errors: bool):
        return

    def display_error_msg(self, msg_string: str):
        self.alteryx_engine.output_message(self.n_tool_id, Sdk.EngineMessageType.error, msg_string)

    def display_info_msg(self, msg_string: str):
        self.alteryx_engine.output_message(self.n_tool_id, Sdk.EngineMessageType.info, msg_string)


class IncomingInterface:
    def __init__(self, parent: AyxPlugin):
        # Default properties
        self.parent: AyxPlugin = parent

        # Custom properties
        self.InInfo: Sdk.RecordInfo = None
        self.OutInfo: Sdk.RecordInfo = None
        self.Creator: Sdk.RecordCreator = None
        self.Copier: Sdk.RecordCopier = None
        self.ScoreField: Sdk.Field = None
        self.Text1Field: Sdk.Field = None
        self.Text2Field: Sdk.Field = None

    def ii_init(self, record_info_in: Sdk.RecordInfo) -> bool:
        self.InInfo = record_info_in
        self.Text1Field = record_info_in.get_field_by_name(self.parent.Text1)
        self.Text2Field = record_info_in.get_field_by_name(self.parent.Text2)
        self.OutInfo = self.InInfo.clone()
        self.ScoreField = self.OutInfo.add_field(self.parent.OutputField, Sdk.FieldType.double, source=self.parent.label)
        self.Creator = self.OutInfo.construct_record_creator()
        self.Copier = Sdk.RecordCopier(self.OutInfo, self.InInfo)

        index = 0
        while index < self.InInfo.num_fields:
            self.Copier.add(index, index)
            index += 1
        self.Copier.done_adding()
        self.parent.Output.init(self.OutInfo)
        return True

    def ii_push_record(self, in_record: Sdk.RecordRef) -> bool:
        self.Creator.reset()
        self.Copier.copy(self.Creator, in_record)
        text1 = self.Text1Field.get_as_string(in_record)
        text2 = self.Text2Field.get_as_string(in_record)
        self.ScoreField.set_from_double(self.Creator, dice_coefficient(text1, text2))
        out_record = self.Creator.finalize_record()
        self.parent.Output.push_record(out_record)
        return True

    def ii_update_progress(self, d_percent: float):
        # Inform the Alteryx engine of the tool's progress.
        self.parent.alteryx_engine.output_tool_progress(self.parent.n_tool_id, d_percent)

    def ii_close(self):
        self.parent.Output.assert_close()
        return


def dice_coefficient(a, b):
    if a is None or b is None: return 0.0
    if not len(a) or not len(b): return 0.0
    """ quick case for true duplicates """
    if a == b: return 1.0
    """ if a != b, and a or b are single chars, then they can't possibly match """
    if len(a) == 1 or len(b) == 1: return 0.0

    """ use python list comprehension, preferred over list.append() """
    a_bigram_list = [a[i:i + 2] for i in range(len(a) - 1)]
    b_bigram_list = [b[i:i + 2] for i in range(len(b) - 1)]

    a_bigram_list.sort()
    b_bigram_list.sort()

    # assignments to save function calls
    lena = len(a_bigram_list)
    lenb = len(b_bigram_list)
    # initialize match counters
    matches = i = j = 0
    while (i < lena and j < lenb):
        if a_bigram_list[i] == b_bigram_list[j]:
            matches += 1
            i += 1
            j += 1
        elif a_bigram_list[i] < b_bigram_list[j]:
            i += 1
        else:
            j += 1

    score = float(2 * matches) / float(lena + lenb)
    return score
